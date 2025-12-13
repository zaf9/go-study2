import { act } from "@testing-library/react";
import {
  getProgress,
  getTopicProgress,
  updateProgress,
  __test__,
} from "@/services/progressService";
import { ProgressStatuses } from "@/lib/progressStatus";

jest.mock("@/lib/api", () => {
  const post = jest.fn();
  const get = jest.fn();
  return {
    __esModule: true,
    default: {
      post,
      get,
    },
  };
});

jest.mock("swr", () => ({
  mutate: jest.fn(),
}));

const mockedApi = jest.requireMock("@/lib/api").default as {
  post: jest.Mock;
  get: jest.Mock;
};

describe("progressService", () => {
  beforeEach(() => {
    jest.useFakeTimers();
    mockedApi.post.mockReset();
    mockedApi.get.mockReset();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("retries updateProgress with backoff", async () => {
    mockedApi.post
      .mockRejectedValueOnce(new Error("network"))
      .mockRejectedValueOnce(new Error("network"))
      .mockResolvedValue({ status: "ok" });

    const promise = updateProgress({ topic: "variables", chapter: "storage" });
    await act(async () => {
      await jest.runAllTimersAsync();
    });
    await promise;

    expect(mockedApi.post).toHaveBeenCalledTimes(3);
  });

  it("returns overall defaults when response missing", async () => {
    mockedApi.get.mockResolvedValue(undefined);
    const snapshot = await getProgress();
    expect(snapshot.overall.totalChapters).toBe(0);
    expect(Array.isArray(snapshot.topics)).toBe(true);
  });

  it("normalizes topic progress response", async () => {
    mockedApi.get.mockResolvedValue({
      topic: {
        id: "variables",
        name: "Variables",
        weight: 25,
        progress: 40,
        completedChapters: 1,
        totalChapters: 4,
      },
      chapters: [{ chapter: "storage", status: ProgressStatuses.InProgress }],
    });
    const detail = await getTopicProgress("variables");
    expect(detail.id).toBe("variables");
    expect(detail.chapters[0].chapter).toBe("storage");
  });

  it("exposes jitter delay helper", () => {
    const delay = __test__.jitterDelay(1);
    expect(delay).toBeGreaterThan(0);
  });
});

