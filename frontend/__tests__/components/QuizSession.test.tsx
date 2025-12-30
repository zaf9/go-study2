/**
 * T063-T064: QuizSession 组件测试
 * 测试测验组件渲染、答题交互和防重复提交逻辑
 * 这里测试 QuizQuestionCard 作为核心测验组件
 */

import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import QuizQuestionCard from "@/components/quiz/QuizQuestionCard";
import { QuizQuestion } from "@/types/quiz";

describe("QuizSession组件 (QuizQuestionCard)", () => {
    const mockSingleQuestion: QuizQuestion = {
        id: 1,
        type: "single",
        difficulty: "easy",
        question: "Which keyword is used to declare a constant in Go?",
        options: ["var", "const", "let", "define"],
        codeSnippet: null,
    };

    const mockMultipleQuestion: QuizQuestion = {
        id: 2,
        type: "multiple",
        difficulty: "medium",
        question: "Which of the following are value types in Go?",
        options: ["int", "string", "slice", "array"],
        codeSnippet: null,
    };

    const mockCodeOutputQuestion: QuizQuestion = {
        id: 3,
        type: "code_output",
        difficulty: "hard",
        question: "What is the output of the following code?",
        options: ["0", "1", "nil", "compile error"],
        codeSnippet: `package main
import "fmt"

func main() {
    var x int
    fmt.Println(x)
}`,
    };

    describe("T063: 测验组件渲染和交互", () => {
        it("应该正确渲染单选题", () => {
            render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={[]}
                    onChange={jest.fn()}
                />,
            );

            // 验证题目文本
            expect(
                screen.getByText(
                    /Which keyword is used to declare a constant in Go\?/i,
                ),
            ).toBeInTheDocument();

            // 验证题型标签
            expect(screen.getByText("单选")).toBeInTheDocument();

            // 验证选项
            expect(screen.getByText(/var/)).toBeInTheDocument();
            expect(screen.getByText(/const/)).toBeInTheDocument();
            expect(screen.getByText(/let/)).toBeInTheDocument();
            expect(screen.getByText(/define/)).toBeInTheDocument();

            // 验证引导语
            expect(
                screen.getByText(/单选题：请选择一个最合适的答案。/),
            ).toBeInTheDocument();
        });

        it("应该正确渲染多选题", () => {
            render(
                <QuizQuestionCard
                    question={mockMultipleQuestion}
                    value={[]}
                    onChange={jest.fn()}
                />,
            );

            // 验证题型标签
            expect(screen.getByText("多选")).toBeInTheDocument();

            // 验证题目文本
            expect(
                screen.getByText(
                    /Which of the following are value types in Go\?/i,
                ),
            ).toBeInTheDocument();

            // 验证多选引导语
            expect(
                screen.getByText(
                    /多选题：请选择一个或多个答案。选中后内容背景将高亮。/,
                ),
            ).toBeInTheDocument();
        });

        it("应该正确渲染代码输出题及代码片段", () => {
            render(
                <QuizQuestionCard
                    question={mockCodeOutputQuestion}
                    value={[]}
                    onChange={jest.fn()}
                />,
            );

            // 验证题型标签
            expect(screen.getByText("代码输出")).toBeInTheDocument();

            // 验证题目文本
            expect(
                screen.getByText(/What is the output of the following code\?/i),
            ).toBeInTheDocument();

            // 验证代码片段存在
            expect(screen.getByText(/package main/)).toBeInTheDocument();
            expect(screen.getByText(/func main\(\)/)).toBeInTheDocument();
        });

        it("单选题选择答案后应该调用onChange", () => {
            const handleChange = jest.fn();
            render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={[]}
                    onChange={handleChange}
                />,
            );

            // 点击选项 "const"
            const constOption = screen.getByText(/const/);
            fireEvent.click(constOption);

            // 验证onChange被调用
            expect(handleChange).toHaveBeenCalledWith(["const"]);
        });

        it("多选题选择多个答案后应该调用onChange", () => {
            const handleChange = jest.fn();
            render(
                <QuizQuestionCard
                    question={mockMultipleQuestion}
                    value={[]}
                    onChange={handleChange}
                />,
            );

            // 选择 "int"
            const intCheckbox = screen.getAllByRole("checkbox")[0];
            fireEvent.click(intCheckbox);

            // 验证onChange被调用，包含 "int"
            expect(handleChange).toHaveBeenCalledWith(["int"]);
        });

        it("应该正确显示已选择的答案", () => {
            render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={["const"]}
                    onChange={jest.fn()}
                />,
            );

            // 找到选中的radio
            const checkedRadio = screen.getByRole("radio", { checked: true });
            expect(checkedRadio).toBeInTheDocument();
        });

        it("多选题应该正确显示多个已选择的答案", () => {
            render(
                <QuizQuestionCard
                    question={mockMultipleQuestion}
                    value={["int", "array"]}
                    onChange={jest.fn()}
                />,
            );

            // 找到所有选中的checkbox
            const checkedCheckboxes = screen.getAllByRole("checkbox", {
                checked: true,
            });
            expect(checkedCheckboxes).toHaveLength(2);
        });

        it("禁用状态下应该无法选择答案", () => {
            const handleChange = jest.fn();
            render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={[]}
                    onChange={handleChange}
                    disabled={true}
                />,
            );

            // 尝试点击选项
            const constOption = screen.getByText(/const/);
            fireEvent.click(constOption);

            // onChange不应该被调用
            expect(handleChange).not.toHaveBeenCalled();
        });
    });

    describe("T064: 防重复提交逻辑", () => {
        it("已禁用状态下不应该响应用户点击", () => {
            const handleChange = jest.fn();
            render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={[]}
                    onChange={handleChange}
                    disabled={true}
                />,
            );

            const varOption = screen.getByText(/var/);
            const constOption = screen.getByText(/const/);

            // 多次点击
            fireEvent.click(varOption);
            fireEvent.click(constOption);
            fireEvent.click(varOption);

            // onChange不应该被调用
            expect(handleChange).not.toHaveBeenCalled();
        });

        it("多选题禁用状态下checkbox应该无法勾选", () => {
            const handleChange = jest.fn();
            render(
                <QuizQuestionCard
                    question={mockMultipleQuestion}
                    value={[]}
                    onChange={handleChange}
                    disabled={true}
                />,
            );

            const checkboxes = screen.getAllByRole("checkbox");

            // 尝试点击多个checkbox
            checkboxes.forEach((checkbox) => {
                fireEvent.click(checkbox);
            });

            // onChange不应该被调用
            expect(handleChange).not.toHaveBeenCalled();
        });

        it("禁用后再启用应该可以正常选择", () => {
            const handleChange = jest.fn();
            const { rerender } = render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={[]}
                    onChange={handleChange}
                    disabled={true}
                />,
            );

            // 禁用状态下点击
            const constOption = screen.getByText(/const/);
            fireEvent.click(constOption);
            expect(handleChange).not.toHaveBeenCalled();

            // 重新渲染为启用状态
            rerender(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={[]}
                    onChange={handleChange}
                    disabled={false}
                />,
            );

            // 再次点击
            fireEvent.click(constOption);
            expect(handleChange).toHaveBeenCalledWith(["const"]);
        });

        it("已选择答案的禁用状态应该保留选择", () => {
            render(
                <QuizQuestionCard
                    question={mockSingleQuestion}
                    value={["const"]}
                    onChange={jest.fn()}
                    disabled={true}
                />,
            );

            // 验证选项仍然是选中的
            const checkedRadio = screen.getByRole("radio", { checked: true });
            expect(checkedRadio).toBeInTheDocument();
        });

        it("禁用状态下多选题已选择的答案应该保留", () => {
            render(
                <QuizQuestionCard
                    question={mockMultipleQuestion}
                    value={["int", "array"]}
                    onChange={jest.fn()}
                    disabled={true}
                />,
            );

            // 验证两个选项仍然是选中的
            const checkedCheckboxes = screen.getAllByRole("checkbox", {
                checked: true,
            });
            expect(checkedCheckboxes).toHaveLength(2);
        });
    });
});
