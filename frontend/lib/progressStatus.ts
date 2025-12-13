import { ProgressStatus } from "@/types/learning";

export const ProgressStatuses = {
  NotStarted: "not_started" as ProgressStatus,
  InProgress: "in_progress" as ProgressStatus,
  Completed: "completed" as ProgressStatus,
  Tested: "tested" as ProgressStatus,
} as const;

export type ProgressStatusKey = keyof typeof ProgressStatuses;
export type ProgressStatusValue = (typeof ProgressStatuses)[ProgressStatusKey];

export default ProgressStatuses;
