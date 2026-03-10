from __future__ import annotations

from dataclasses import dataclass, asdict
from datetime import datetime, timedelta

from .api import DayMetrics


# ── Result types ─────────────────────────────────────────────────────────────


@dataclass
class Summary:
    days_count: int = 0
    total_active_users: int = 0
    total_engaged_users: int = 0
    avg_active_users: float = 0.0
    avg_engaged_users: float = 0.0
    max_active_users: int = 0
    total_suggestions: int = 0
    total_acceptances: int = 0
    total_lines_suggested: int = 0
    total_lines_accepted: int = 0
    acceptance_rate: float = 0.0
    total_chats: int = 0
    total_chat_insertions: int = 0
    total_chat_copies: int = 0
    total_dotcom_chats: int = 0
    total_pr_summaries: int = 0


@dataclass
class LanguageStat:
    name: str = ""
    suggestions: int = 0
    acceptances: int = 0
    lines_suggested: int = 0
    lines_accepted: int = 0
    acceptance_rate: float = 0.0


@dataclass
class EditorStat:
    name: str = ""
    engaged_users: int = 0
    suggestions: int = 0
    acceptances: int = 0
    acceptance_rate: float = 0.0
    chats: int = 0


@dataclass
class DaySummary:
    date: str = ""
    active_users: int = 0
    engaged_users: int = 0
    suggestions: int = 0
    acceptances: int = 0
    acceptance_rate: float = 0.0
    lines_accepted: int = 0
    chats: int = 0


@dataclass
class WeekSummary:
    label: str = ""
    start_date: str = ""
    end_date: str = ""
    days: int = 0
    avg_active_users: float = 0.0
    total_suggestions: int = 0
    total_acceptances: int = 0
    acceptance_rate: float = 0.0
    total_chats: int = 0


def to_dict(obj):
    """Convert a dataclass (or list of dataclasses) to plain dicts."""
    if isinstance(obj, list):
        return [asdict(item) for item in obj]
    return asdict(obj)


# ── Summarize ────────────────────────────────────────────────────────────────


def summarize(days: list[DayMetrics]) -> Summary:
    s = Summary(days_count=len(days))
    if not days:
        return s

    for d in days:
        s.total_active_users += d.total_active_users
        s.total_engaged_users += d.total_engaged_users
        if d.total_active_users > s.max_active_users:
            s.max_active_users = d.total_active_users

        _add_code_completions(s, d)
        _add_ide_chat(s, d)
        _add_dotcom_chat(s, d)
        _add_pr_summaries(s, d)

    s.avg_active_users = s.total_active_users / len(days)
    s.avg_engaged_users = s.total_engaged_users / len(days)
    if s.total_suggestions > 0:
        s.acceptance_rate = s.total_acceptances / s.total_suggestions * 100

    return s


def _add_code_completions(s: Summary, d: DayMetrics):
    cc = d.copilot_ide_code_completions
    if cc is None:
        return
    for editor in cc.editors:
        for model in editor.models:
            for lang in model.languages:
                s.total_suggestions += lang.total_code_suggestions
                s.total_acceptances += lang.total_code_acceptances
                s.total_lines_suggested += lang.total_code_lines_suggested
                s.total_lines_accepted += lang.total_code_lines_accepted


def _add_ide_chat(s: Summary, d: DayMetrics):
    chat = d.copilot_ide_chat
    if chat is None:
        return
    for editor in chat.editors:
        for model in editor.models:
            s.total_chats += model.total_chats
            s.total_chat_insertions += model.total_chat_insertion_events
            s.total_chat_copies += model.total_chat_copy_events


def _add_dotcom_chat(s: Summary, d: DayMetrics):
    chat = d.copilot_dotcom_chat
    if chat is None:
        return
    for model in chat.models:
        s.total_dotcom_chats += model.total_chats


def _add_pr_summaries(s: Summary, d: DayMetrics):
    pr = d.copilot_dotcom_pull_requests
    if pr is None:
        return
    for repo in pr.repositories:
        for model in repo.models:
            s.total_pr_summaries += model.total_pr_summaries


# ── By Language ──────────────────────────────────────────────────────────────


def by_language(days: list[DayMetrics]) -> list[LanguageStat]:
    lang_map: dict[str, LanguageStat] = {}

    for d in days:
        cc = d.copilot_ide_code_completions
        if cc is None:
            continue
        for editor in cc.editors:
            for model in editor.models:
                for lang in model.languages:
                    ls = lang_map.get(lang.name)
                    if ls is None:
                        ls = LanguageStat(name=lang.name)
                        lang_map[lang.name] = ls
                    ls.suggestions += lang.total_code_suggestions
                    ls.acceptances += lang.total_code_acceptances
                    ls.lines_suggested += lang.total_code_lines_suggested
                    ls.lines_accepted += lang.total_code_lines_accepted

    result = list(lang_map.values())
    for ls in result:
        if ls.suggestions > 0:
            ls.acceptance_rate = ls.acceptances / ls.suggestions * 100

    result.sort(key=lambda x: x.acceptances, reverse=True)
    return result


# ── By Editor ────────────────────────────────────────────────────────────────


def by_editor(days: list[DayMetrics]) -> list[EditorStat]:
    editor_map: dict[str, EditorStat] = {}

    for d in days:
        cc = d.copilot_ide_code_completions
        if cc is not None:
            for editor in cc.editors:
                es = editor_map.get(editor.name)
                if es is None:
                    es = EditorStat(name=editor.name)
                    editor_map[editor.name] = es
                for model in editor.models:
                    for lang in model.languages:
                        es.suggestions += lang.total_code_suggestions
                        es.acceptances += lang.total_code_acceptances

        chat = d.copilot_ide_chat
        if chat is not None:
            for editor in chat.editors:
                es = editor_map.get(editor.name)
                if es is None:
                    es = EditorStat(name=editor.name)
                    editor_map[editor.name] = es
                for model in editor.models:
                    es.chats += model.total_chats

    # Track max engaged users per editor (max across days, not sum)
    for d in days:
        cc = d.copilot_ide_code_completions
        if cc is not None:
            for editor in cc.editors:
                es = editor_map.get(editor.name)
                if es and editor.total_engaged_users > es.engaged_users:
                    es.engaged_users = editor.total_engaged_users

    result = list(editor_map.values())
    for es in result:
        if es.suggestions > 0:
            es.acceptance_rate = es.acceptances / es.suggestions * 100

    result.sort(key=lambda x: x.suggestions, reverse=True)
    return result


# ── Daily Breakdown ──────────────────────────────────────────────────────────


def daily_breakdown(days: list[DayMetrics]) -> list[DaySummary]:
    result: list[DaySummary] = []

    for d in days:
        ds = DaySummary(
            date=d.date,
            active_users=d.total_active_users,
            engaged_users=d.total_engaged_users,
        )

        cc = d.copilot_ide_code_completions
        if cc is not None:
            for editor in cc.editors:
                for model in editor.models:
                    for lang in model.languages:
                        ds.suggestions += lang.total_code_suggestions
                        ds.acceptances += lang.total_code_acceptances
                        ds.lines_accepted += lang.total_code_lines_accepted

        chat = d.copilot_ide_chat
        if chat is not None:
            for editor in chat.editors:
                for model in editor.models:
                    ds.chats += model.total_chats

        if ds.suggestions > 0:
            ds.acceptance_rate = ds.acceptances / ds.suggestions * 100

        result.append(ds)

    return result


# ── Week-over-Week ───────────────────────────────────────────────────────────


def week_over_week(
    days: list[DayMetrics],
) -> tuple[WeekSummary, WeekSummary]:
    now = datetime.now()

    weekday = now.weekday()  # Monday=0
    this_monday = now - timedelta(days=weekday)
    last_monday = this_monday - timedelta(days=7)
    last_sunday = this_monday - timedelta(days=1)

    this_mon_str = this_monday.strftime("%Y-%m-%d")
    last_mon_str = last_monday.strftime("%Y-%m-%d")
    last_sun_str = last_sunday.strftime("%Y-%m-%d")
    today_str = now.strftime("%Y-%m-%d")

    this_days = [d for d in days if this_mon_str <= d.date <= today_str]
    last_days = [d for d in days if last_mon_str <= d.date <= last_sun_str]

    this_week = _week_summary("This week", this_mon_str, today_str, this_days)
    last_week = _week_summary("Last week", last_mon_str, last_sun_str, last_days)

    return this_week, last_week


def _week_summary(
    label: str, start: str, end: str, days: list[DayMetrics]
) -> WeekSummary:
    s = summarize(days)
    return WeekSummary(
        label=label,
        start_date=start,
        end_date=end,
        days=len(days),
        avg_active_users=s.avg_active_users,
        total_suggestions=s.total_suggestions,
        total_acceptances=s.total_acceptances,
        acceptance_rate=s.acceptance_rate,
        total_chats=s.total_chats + s.total_dotcom_chats,
    )
