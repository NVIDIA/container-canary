---
name: issue-writer
description: Draft and revise concise, human-focused GitHub issues for pytest-kind-ng. Use when asked to write, rewrite, review, or prepare a bug report, enhancement request, maintenance issue, or follow-up issue for this repository, including turning code-review findings into issues.
---

# Issue Writer

Write issues that let a maintainer understand the problem, reproduce it when
applicable, and decide whether it matters without reading an implementation plan.

## Workflow

1. Inspect the relevant code, tests, documentation, and repository guidance. Do
   not invent behavior, impact, versions, or reproduction results.
2. Reduce the request to one problem. Split unrelated problems into separate
   issues.
3. Remove proposed fixes and implementation details. Preserve constraints and
   observable desired behavior only when they clarify the problem.
4. For a bug, build the smallest complete reproduction and verify that it still
   demonstrates the reported behavior. For a non-bug change, state the concrete
   limitation and why resolving it is important.
5. Draft only the sections that add information, then edit for brevity.
6. Return a draft unless the user explicitly asks to create or submit the issue.

## Write for Maintainers

- Lead with the observed problem and its impact.
- Use a specific, problem-focused title. Describe what fails or is difficult,
  not a technology or refactoring to adopt.
- Keep paragraphs short and omit background that does not change understanding
  or reproduction.
- Use plain language for humans. Do not include agent instructions, exhaustive
  task lists, review history, or commentary about how the issue was produced.
- Distinguish verified facts from uncertainty. Do not speculate about a root
  cause.
- Do not add labels, priority, milestones, or assignees unless the user provides
  them.

Prefer titles such as:

> Port forwarding restarts after it becomes ready

Avoid solution-shaped titles such as:

> Refactor port forwarding to use a new readiness mechanism

## Reproduce Bugs

Make reproductions:

- **Minimal:** Remove every line, step, fixture, dependency, and input that is not
  required to trigger the behavior.
- **Complete:** Include everything a maintainer needs to copy, run, and observe
  the problem.
- **Verifiable:** Run the final reduced example and report the command and result.
  Never claim that an untested example reproduces the problem.

Prefer synthetic inline configuration or data over external repositories,
archives, logs, or private application data. Include only relevant environment
details, such as Python, pytest, kind, kubectl, OS, and architecture versions.
Use language-labelled fenced code blocks. Include the complete traceback or
command output when it is relevant; place unusually long output in a
`<details>` block rather than truncating the useful frames.

Present Python examples in `python`-labelled code blocks. Do not wrap Python
code in a shell heredoc such as `python3 <<'PY' ... PY`; readers can infer that
the example should be run with Python. Use a `bash` block only for shell
commands.

If a minimal reproduction is relevant but unavailable, state that limitation
plainly instead of fabricating one.

## Structure the Draft

Adapt the structure to the issue. Do not emit empty headings.

For a bug, usually include:

- A one-paragraph summary
- A minimal reproduction with the exact command
- Actual behavior, including the complete relevant error
- Expected behavior
- Relevant environment details

For an enhancement or maintenance issue, usually include:

- The current limitation or maintenance problem
- A concrete example or affected workflow
- Why it matters to users or maintainers
- Clear boundaries when needed to keep the issue focused

Do not add a “Proposed solution” section, name files to edit, prescribe an
architecture, or provide implementation steps. State the outcome that should be
possible only when it is necessary to explain the problem.

## Final Check

Before returning the draft, verify:

- The title describes one problem rather than a solution.
- The opening explains what is wrong and why it matters.
- A bug reproduction is minimal, complete, and actually verified.
- A non-bug issue gives concrete justification.
- Expected behavior is observable without prescribing implementation.
- Every paragraph earns the maintainer's attention.
