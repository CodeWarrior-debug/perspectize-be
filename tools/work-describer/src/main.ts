/**
 * Work Describer / Prompt Builder
 *
 * Fully deterministic: builds a structured work-request prompt from form
 * fields using plain string templates. No network calls, no AI/tokens.
 */

type WorkingLocation = "main-worktree" | "available-worktree";
type TestingApproach = "codebase-then-judgment" | "stop-and-ask";

interface FormState {
  taskName: string;
  taskType: string;
  intent: string;
  links: string;
  repo: string;
  startingBranch: string;
  workingLocation: WorkingLocation;
  prTargetBranch: string;
  testingApproach: TestingApproach;
  testingNotes: string;
  ignoreInstructions: string;
}

const STORAGE_KEY = "work-describer:v1";
const REPO_LIST_KEY = "work-describer:repos:v1";
const BRANCH_LIST_KEY = "work-describer:branches:v1";

const DEFAULT_REPOS = ["perspectize"];
const DEFAULT_BRANCHES = ["main"];

const TASK_TYPE_LABELS: Record<string, string> = {
  feature: "Feature",
  bugfix: "Bugfix",
  chore: "Chore",
  docs: "Docs",
  refactor: "Refactor",
  test: "Test",
  other: "Other",
};

const WORKING_LOCATION_LABELS: Record<WorkingLocation, string> = {
  "main-worktree": "Main worktree",
  "available-worktree": "Available non-main worktree",
};

function $<T extends HTMLElement>(id: string): T {
  const el = document.getElementById(id);
  if (!el) throw new Error(`Missing element #${id}`);
  return el as T;
}

// ---------------------------------------------------------------------------
// Persisted repo / branch lists (so "latest repo" can be set as a default)
// ---------------------------------------------------------------------------

function loadList(key: string, fallback: string[]): string[] {
  try {
    const raw = localStorage.getItem(key);
    if (!raw) return fallback.slice();
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed) && parsed.every((x) => typeof x === "string") && parsed.length > 0) {
      return parsed;
    }
  } catch {
    /* ignore malformed storage */
  }
  return fallback.slice();
}

function saveList(key: string, list: string[]): void {
  try {
    localStorage.setItem(key, JSON.stringify(list));
  } catch {
    /* localStorage unavailable — non-fatal */
  }
}

function populateDatalist(datalistId: string, values: string[]): void {
  const dl = $(datalistId);
  dl.innerHTML = "";
  for (const v of values) {
    const opt = document.createElement("option");
    opt.value = v;
    dl.appendChild(opt);
  }
}

function addToFront(list: string[], value: string): string[] {
  const trimmed = value.trim();
  if (!trimmed) return list;
  const rest = list.filter((v) => v !== trimmed);
  return [trimmed, ...rest].slice(0, 20);
}

// ---------------------------------------------------------------------------
// Form <-> state
// ---------------------------------------------------------------------------

function readForm(): FormState {
  return {
    taskName: $<HTMLInputElement>("taskName").value.trim(),
    taskType: $<HTMLSelectElement>("taskType").value,
    intent: $<HTMLTextAreaElement>("intent").value.trim(),
    links: $<HTMLTextAreaElement>("links").value.trim(),
    repo: $<HTMLInputElement>("repo").value.trim(),
    startingBranch: $<HTMLInputElement>("startingBranch").value.trim(),
    workingLocation: $<HTMLSelectElement>("workingLocation").value as WorkingLocation,
    prTargetBranch: $<HTMLInputElement>("prTargetBranch").value.trim(),
    testingApproach: (document.querySelector<HTMLInputElement>(
      'input[name="testingApproach"]:checked',
    )?.value ?? "codebase-then-judgment") as TestingApproach,
    testingNotes: $<HTMLTextAreaElement>("testingNotes").value.trim(),
    ignoreInstructions: $<HTMLTextAreaElement>("ignoreInstructions").value.trim(),
  };
}

function writeForm(state: Partial<FormState>): void {
  if (state.taskName !== undefined) $<HTMLInputElement>("taskName").value = state.taskName;
  if (state.taskType !== undefined) $<HTMLSelectElement>("taskType").value = state.taskType;
  if (state.intent !== undefined) $<HTMLTextAreaElement>("intent").value = state.intent;
  if (state.links !== undefined) $<HTMLTextAreaElement>("links").value = state.links;
  if (state.repo !== undefined) $<HTMLInputElement>("repo").value = state.repo;
  if (state.startingBranch !== undefined)
    $<HTMLInputElement>("startingBranch").value = state.startingBranch;
  if (state.workingLocation !== undefined)
    $<HTMLSelectElement>("workingLocation").value = state.workingLocation;
  if (state.prTargetBranch !== undefined)
    $<HTMLInputElement>("prTargetBranch").value = state.prTargetBranch;
  if (state.testingApproach !== undefined) {
    const radio = document.querySelector<HTMLInputElement>(
      `input[name="testingApproach"][value="${state.testingApproach}"]`,
    );
    if (radio) radio.checked = true;
  }
  if (state.testingNotes !== undefined)
    $<HTMLTextAreaElement>("testingNotes").value = state.testingNotes;
  if (state.ignoreInstructions !== undefined)
    $<HTMLTextAreaElement>("ignoreInstructions").value = state.ignoreInstructions;
}

function loadPersistedState(): Partial<FormState> | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as Partial<FormState>;
  } catch {
    return null;
  }
}

function persistState(state: FormState): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    /* non-fatal */
  }
}

// ---------------------------------------------------------------------------
// Deterministic prompt rendering
// ---------------------------------------------------------------------------

function parseLines(text: string): string[] {
  return text
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l.length > 0);
}

function renderMarkdown(state: FormState): string {
  const taskTypeLabel = TASK_TYPE_LABELS[state.taskType] ?? state.taskType;
  const title = state.taskName || "(untitled task)";
  const lines: string[] = [];

  lines.push(`# ${title}`);
  lines.push("");
  lines.push(`**Task type:** ${taskTypeLabel}`);
  lines.push("");

  lines.push("## Intent");
  lines.push("");
  lines.push(state.intent || "_(not provided)_");
  lines.push("");

  const linkLines = parseLines(state.links);
  if (linkLines.length > 0) {
    lines.push("## Links & attachments");
    lines.push("");
    for (const l of linkLines) {
      lines.push(`- ${l}`);
    }
    lines.push("");
  }

  lines.push("## Environment");
  lines.push("");
  lines.push(`- **Repo:** ${state.repo || "_(not specified)_"}`);
  lines.push(`- **Starting branch:** ${state.startingBranch || "main"}`);
  lines.push(
    `- **Working location:** ${WORKING_LOCATION_LABELS[state.workingLocation] ?? state.workingLocation}`,
  );
  lines.push(`- **PR target branch:** ${state.prTargetBranch || "main"}`);
  lines.push("");

  lines.push("## Testing approach");
  lines.push("");
  if (state.testingApproach === "stop-and-ask") {
    lines.push("Stop and ask if really unsure — do not guess at a testing approach.");
  } else {
    lines.push(
      "Follow codebase instructions first, then established patterns, then judgment. Stop and ask if really unsure.",
    );
  }
  if (state.testingNotes) {
    lines.push("");
    lines.push(`Additional notes: ${state.testingNotes}`);
  }
  lines.push("");

  if (state.ignoreInstructions) {
    lines.push("## Ignore for this session only");
    lines.push("");
    lines.push(
      "The following codebase instructions/approaches should be ignored for this specific session only:",
    );
    lines.push("");
    for (const l of parseLines(state.ignoreInstructions)) {
      lines.push(`- ${l}`);
    }
    lines.push("");
  }

  return lines.join("\n").trimEnd() + "\n";
}

/** Minimal, deterministic markdown -> HTML for the rich-text preview pane. */
function markdownToHtml(md: string): string {
  const escape = (s: string) =>
    s
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;");

  const linkify = (s: string) =>
    s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>');

  const inline = (s: string) => {
    let out = escape(s);
    out = out.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
    out = linkify(out);
    return out;
  };

  const rawLines = md.split("\n");
  const html: string[] = [];
  let inList = false;

  const closeList = () => {
    if (inList) {
      html.push("</ul>");
      inList = false;
    }
  };

  for (const line of rawLines) {
    if (line.startsWith("# ")) {
      closeList();
      html.push(`<h1>${inline(line.slice(2))}</h1>`);
    } else if (line.startsWith("## ")) {
      closeList();
      html.push(`<h2>${inline(line.slice(3))}</h2>`);
    } else if (line.startsWith("- ")) {
      if (!inList) {
        html.push("<ul>");
        inList = true;
      }
      html.push(`<li>${inline(line.slice(2))}</li>`);
    } else if (line.trim() === "") {
      closeList();
    } else {
      closeList();
      html.push(`<p>${inline(line)}</p>`);
    }
  }
  closeList();
  return html.join("\n");
}

// ---------------------------------------------------------------------------
// Wiring
// ---------------------------------------------------------------------------

function update(): void {
  const state = readForm();
  persistState(state);
  const md = renderMarkdown(state);
  $<HTMLTextAreaElement>("rawOutput").value = md;
  $<HTMLDivElement>("richPreview").innerHTML = markdownToHtml(md);
}

function slugify(name: string): string {
  return (
    name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/^-+|-+$/g, "") || "work-request"
  );
}

function download(): void {
  const state = readForm();
  const md = renderMarkdown(state);
  const blob = new Blob([md], { type: "text/markdown;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `${slugify(state.taskName)}.md`;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(url);
}

async function copyToClipboard(): Promise<void> {
  const text = $<HTMLTextAreaElement>("rawOutput").value;
  const btn = $<HTMLButtonElement>("copyBtn");
  const original = btn.textContent;
  try {
    await navigator.clipboard.writeText(text);
    btn.textContent = "Copied!";
  } catch {
    // Fallback for environments without Clipboard API permission.
    const ta = $<HTMLTextAreaElement>("rawOutput");
    ta.focus();
    ta.select();
    document.execCommand("copy");
    btn.textContent = "Copied!";
  }
  window.setTimeout(() => {
    btn.textContent = original;
  }, 1500);
}

function init(): void {
  const repos = loadList(REPO_LIST_KEY, DEFAULT_REPOS);
  const branches = loadList(BRANCH_LIST_KEY, DEFAULT_BRANCHES);
  populateDatalist("repoList", repos);
  populateDatalist("branchList", branches);

  const persisted = loadPersistedState();
  writeForm({
    repo: repos[0],
    startingBranch: branches[0],
    prTargetBranch: "main",
    workingLocation: "main-worktree",
    testingApproach: "codebase-then-judgment",
    ...persisted,
  });

  const formEl = $<HTMLFormElement>("describerForm");
  formEl.addEventListener("input", update);
  formEl.addEventListener("change", update);

  $<HTMLButtonElement>("copyBtn").addEventListener("click", () => void copyToClipboard());
  $<HTMLButtonElement>("downloadBtn").addEventListener("click", download);

  $<HTMLButtonElement>("resetBtn").addEventListener("click", () => {
    formEl.reset();
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch {
      /* non-fatal */
    }
    writeForm({
      repo: repos[0],
      startingBranch: branches[0],
      prTargetBranch: "main",
      workingLocation: "main-worktree",
      testingApproach: "codebase-then-judgment",
    });
    update();
  });

  $<HTMLButtonElement>("rememberRepoBtn").addEventListener("click", () => {
    const value = $<HTMLInputElement>("repo").value;
    const updated = addToFront(repos, value);
    repos.length = 0;
    repos.push(...updated);
    saveList(REPO_LIST_KEY, repos);
    populateDatalist("repoList", repos);
  });

  $<HTMLButtonElement>("rememberBranchBtn").addEventListener("click", () => {
    const value = $<HTMLInputElement>("startingBranch").value;
    const updated = addToFront(branches, value);
    branches.length = 0;
    branches.push(...updated);
    saveList(BRANCH_LIST_KEY, branches);
    populateDatalist("branchList", branches);
  });

  update();
}

document.addEventListener("DOMContentLoaded", init);
