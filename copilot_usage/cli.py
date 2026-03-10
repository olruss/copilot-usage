from __future__ import annotations

from datetime import datetime, timedelta

import click

from .api import Client
from .auth import get_token
from .analysis import (
    summarize,
    by_language,
    by_editor,
    daily_breakdown,
    week_over_week,
    to_dict,
)
from .display import render_table, render_kv, render_bar_chart, render_json


# ── Helpers ──────────────────────────────────────────────────────────────────


def resolve_date_range(
    period: str | None, since: str | None, until: str | None
) -> tuple[str, str]:
    if since or until:
        return since or "", until or ""

    now = datetime.now()
    today = now.strftime("%Y-%m-%d")

    match period:
        case "today":
            return today, today
        case "yesterday":
            y = (now - timedelta(days=1)).strftime("%Y-%m-%d")
            return y, y
        case "week":
            return (now - timedelta(days=7)).strftime("%Y-%m-%d"), today
        case "month":
            first = now.replace(day=1).strftime("%Y-%m-%d")
            return first, today
        case "last-month":
            first_this = now.replace(day=1)
            last_prev = first_this - timedelta(days=1)
            first_prev = last_prev.replace(day=1).strftime("%Y-%m-%d")
            return first_prev, last_prev.strftime("%Y-%m-%d")
        case "year":
            first_of_year = datetime(now.year, 1, 1)
            days_since = (now - first_of_year).days
            if days_since > 100:
                return (now - timedelta(days=100)).strftime("%Y-%m-%d"), today
            return first_of_year.strftime("%Y-%m-%d"), today
        case _:
            return (now - timedelta(days=28)).strftime("%Y-%m-%d"), today


def format_num(n: int) -> str:
    if n < 1000:
        return str(n)
    if n < 1_000_000:
        return f"{n / 1000:.1f}K"
    return f"{n / 1_000_000:.1f}M"


def format_delta(d: int) -> str:
    if d > 0:
        return f"[green]+{format_num(d)}[/]"
    if d < 0:
        return f"[red]{format_num(d)}[/]"
    return "0"


# ── CLI ──────────────────────────────────────────────────────────────────────


@click.group()
@click.option("--org", "-o", required=True, help="GitHub organization name")
@click.option("--token", default=None, help="GitHub token (default: gh auth token)")
@click.option("--json", "json_out", is_flag=True, help="Output as JSON")
@click.option("--since", default=None, help="Start date (YYYY-MM-DD)")
@click.option("--until", default=None, help="End date (YYYY-MM-DD)")
@click.option(
    "--period",
    "-p",
    default=None,
    help="Named period: today, yesterday, week, month, last-month, year",
)
@click.pass_context
def cli(ctx, org, token, json_out, since, until, period):
    """GitHub Copilot usage analytics for your organization."""
    ctx.ensure_object(dict)
    ctx.obj["org"] = org
    ctx.obj["token"] = token
    ctx.obj["json_out"] = json_out
    ctx.obj["since"] = since
    ctx.obj["until"] = until
    ctx.obj["period"] = period


def fetch_metrics(ctx) -> list:
    tok = get_token(ctx.obj["token"])
    s, u = resolve_date_range(ctx.obj["period"], ctx.obj["since"], ctx.obj["until"])
    client = Client(tok)
    return client.fetch_metrics(ctx.obj["org"], s, u)


# ── summary ──────────────────────────────────────────────────────────────────


@cli.command()
@click.pass_context
def summary(ctx):
    """Show usage summary for a period."""
    days = fetch_metrics(ctx)
    if not days:
        click.echo("No data available for the selected period.")
        return

    s = summarize(days)

    if ctx.obj["json_out"]:
        render_json(
            {
                "period": {"from": days[0].date, "to": days[-1].date},
                "days": s.days_count,
                "avg_active_users": s.avg_active_users,
                "avg_engaged_users": s.avg_engaged_users,
                "max_active_users": s.max_active_users,
                "total_suggestions": s.total_suggestions,
                "total_acceptances": s.total_acceptances,
                "acceptance_rate": s.acceptance_rate,
                "total_lines_suggested": s.total_lines_suggested,
                "total_lines_accepted": s.total_lines_accepted,
                "total_chats": s.total_chats,
                "total_dotcom_chats": s.total_dotcom_chats,
                "total_chat_insertions": s.total_chat_insertions,
                "total_chat_copies": s.total_chat_copies,
                "total_pr_summaries": s.total_pr_summaries,
            }
        )
        return

    period_label = (
        f"{days[0].date} to {days[-1].date} ({s.days_count} days)"
    )

    render_kv(
        "Copilot Usage Summary",
        [
            ("Period", period_label),
            ("Avg Active Users", f"{s.avg_active_users:.1f}"),
            ("Avg Engaged Users", f"{s.avg_engaged_users:.1f}"),
            ("Max Active Users", str(s.max_active_users)),
            ("Suggestions", format_num(s.total_suggestions)),
            ("Acceptances", format_num(s.total_acceptances)),
            ("Acceptance Rate", f"{s.acceptance_rate:.1f}%"),
            ("Lines Suggested", format_num(s.total_lines_suggested)),
            ("Lines Accepted", format_num(s.total_lines_accepted)),
            ("IDE Chats", format_num(s.total_chats)),
            ("Dotcom Chats", format_num(s.total_dotcom_chats)),
            ("Chat Insertions", format_num(s.total_chat_insertions)),
            ("Chat Copies", format_num(s.total_chat_copies)),
            ("PR Summaries", format_num(s.total_pr_summaries)),
        ],
    )


# ── daily ────────────────────────────────────────────────────────────────────


@cli.command()
@click.pass_context
def daily(ctx):
    """Show day-by-day breakdown."""
    days = fetch_metrics(ctx)
    if not days:
        click.echo("No data available for the selected period.")
        return

    breakdown = daily_breakdown(days)

    if ctx.obj["json_out"]:
        render_json(to_dict(breakdown))
        return

    rows = [
        [
            d.date,
            str(d.active_users),
            str(d.engaged_users),
            format_num(d.suggestions),
            format_num(d.acceptances),
            f"{d.acceptance_rate:.1f}%",
            format_num(d.lines_accepted),
            str(d.chats),
        ]
        for d in breakdown
    ]

    render_table(
        "Daily Copilot Usage",
        ["Date", "Active", "Engaged", "Suggestions", "Acceptances", "Rate", "Lines", "Chats"],
        rows,
    )


# ── languages ────────────────────────────────────────────────────────────────


@cli.command(name="languages")
@click.pass_context
def languages_cmd(ctx):
    """Show language breakdown."""
    days = fetch_metrics(ctx)
    if not days:
        click.echo("No data available for the selected period.")
        return

    langs = by_language(days)

    if ctx.obj["json_out"]:
        render_json(to_dict(langs))
        return

    rows = [
        [
            l.name,
            format_num(l.suggestions),
            format_num(l.acceptances),
            f"{l.acceptance_rate:.1f}%",
            format_num(l.lines_accepted),
        ]
        for l in langs
    ]

    render_table(
        f"Language Breakdown ({days[0].date} to {days[-1].date})",
        ["Language", "Suggestions", "Acceptances", "Rate", "Lines Accepted"],
        rows,
    )

    bar_items = [
        (l.name, l.acceptances, f"({l.acceptance_rate:.1f}% rate)") for l in langs
    ]
    render_bar_chart("Acceptances by Language", bar_items)


# alias
cli.add_command(languages_cmd, name="langs")


# ── editors ──────────────────────────────────────────────────────────────────


@cli.command()
@click.pass_context
def editors(ctx):
    """Show editor breakdown."""
    days = fetch_metrics(ctx)
    if not days:
        click.echo("No data available for the selected period.")
        return

    eds = by_editor(days)

    if ctx.obj["json_out"]:
        render_json(to_dict(eds))
        return

    rows = [
        [
            e.name,
            str(e.engaged_users),
            format_num(e.suggestions),
            format_num(e.acceptances),
            f"{e.acceptance_rate:.1f}%",
            str(e.chats),
        ]
        for e in eds
    ]

    render_table(
        f"Editor Breakdown ({days[0].date} to {days[-1].date})",
        ["Editor", "Users", "Suggestions", "Acceptances", "Rate", "Chats"],
        rows,
    )

    bar_items = [
        (e.name, e.suggestions, f"({e.engaged_users} users)") for e in eds
    ]
    render_bar_chart("Suggestions by Editor", bar_items)


# ── trends ───────────────────────────────────────────────────────────────────


@cli.command()
@click.pass_context
def trends(ctx):
    """Show week-over-week trends."""
    tok = get_token(ctx.obj["token"])

    now = datetime.now()
    s = (now - timedelta(days=21)).strftime("%Y-%m-%d")
    u = now.strftime("%Y-%m-%d")

    client = Client(tok)
    days = client.fetch_metrics(ctx.obj["org"], s, u)

    if not days:
        click.echo("No data available.")
        return

    this_week, last_week = week_over_week(days)

    if ctx.obj["json_out"]:
        render_json(
            {
                "this_week": to_dict(this_week),
                "last_week": to_dict(last_week),
            }
        )
        return

    def trend_row(label: str, last: int, this: int) -> list[str]:
        delta = this - last
        return [label, format_num(last), format_num(this), format_delta(delta)]

    def trend_row_pct(label: str, last: float, this: float) -> list[str]:
        delta = this - last
        sign = "+" if delta >= 0 else ""
        return [
            label,
            f"{last:.1f}%",
            f"{this:.1f}%",
            f"{sign}{delta:.1f}%",
        ]

    rows = [
        trend_row("Days with data", last_week.days, this_week.days),
        trend_row(
            "Avg Active Users",
            int(last_week.avg_active_users),
            int(this_week.avg_active_users),
        ),
        trend_row("Suggestions", last_week.total_suggestions, this_week.total_suggestions),
        trend_row("Acceptances", last_week.total_acceptances, this_week.total_acceptances),
        trend_row_pct("Acceptance Rate", last_week.acceptance_rate, this_week.acceptance_rate),
        trend_row("Chats", last_week.total_chats, this_week.total_chats),
    ]

    render_table("Week-over-Week Trends", ["Metric", "Last Week", "This Week", "Delta"], rows)

    click.echo(f"  Last week: {last_week.start_date} to {last_week.end_date}")
    click.echo(f"  This week: {this_week.start_date} to {this_week.end_date}")
    click.echo()
