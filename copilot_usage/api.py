from __future__ import annotations

from dataclasses import dataclass, field

import requests

BASE_URL = "https://api.github.com"
API_VERSION = "2022-11-28"
PER_PAGE = 100


# ── Data types mirroring the GitHub Copilot Metrics API response ─────────────


@dataclass
class IDECodeCompletionsLanguage:
    name: str = ""
    total_engaged_users: int = 0
    total_code_suggestions: int = 0
    total_code_acceptances: int = 0
    total_code_lines_suggested: int = 0
    total_code_lines_accepted: int = 0


@dataclass
class IDECodeCompletionsModel:
    name: str = ""
    is_custom_model: bool = False
    total_engaged_users: int = 0
    languages: list[IDECodeCompletionsLanguage] = field(default_factory=list)


@dataclass
class IDECodeCompletionsEditor:
    name: str = ""
    total_engaged_users: int = 0
    models: list[IDECodeCompletionsModel] = field(default_factory=list)


@dataclass
class IDECodeCompletions:
    total_engaged_users: int = 0
    editors: list[IDECodeCompletionsEditor] = field(default_factory=list)


@dataclass
class IDEChatModel:
    name: str = ""
    is_custom_model: bool = False
    total_engaged_users: int = 0
    total_chats: int = 0
    total_chat_insertion_events: int = 0
    total_chat_copy_events: int = 0


@dataclass
class IDEChatEditor:
    name: str = ""
    total_engaged_users: int = 0
    models: list[IDEChatModel] = field(default_factory=list)


@dataclass
class IDEChat:
    total_engaged_users: int = 0
    editors: list[IDEChatEditor] = field(default_factory=list)


@dataclass
class DotcomChatModel:
    name: str = ""
    is_custom_model: bool = False
    total_engaged_users: int = 0
    total_chats: int = 0


@dataclass
class DotcomChat:
    total_engaged_users: int = 0
    models: list[DotcomChatModel] = field(default_factory=list)


@dataclass
class DotcomPullRequestsModel:
    name: str = ""
    is_custom_model: bool = False
    total_engaged_users: int = 0
    total_pr_summaries: int = 0


@dataclass
class DotcomPullRequestsRepo:
    name: str = ""
    total_engaged_users: int = 0
    models: list[DotcomPullRequestsModel] = field(default_factory=list)


@dataclass
class DotcomPullRequests:
    total_engaged_users: int = 0
    repositories: list[DotcomPullRequestsRepo] = field(default_factory=list)


@dataclass
class DayMetrics:
    date: str = ""
    total_active_users: int = 0
    total_engaged_users: int = 0
    copilot_ide_code_completions: IDECodeCompletions | None = None
    copilot_ide_chat: IDEChat | None = None
    copilot_dotcom_chat: DotcomChat | None = None
    copilot_dotcom_pull_requests: DotcomPullRequests | None = None


# ── Recursive dict→dataclass helper ─────────────────────────────────────────


def _parse(cls, data):
    """Recursively convert a JSON dict into a dataclass instance."""
    if data is None:
        return None
    if not isinstance(data, dict):
        return data

    kwargs = {}
    for f in cls.__dataclass_fields__.values():
        raw = data.get(f.name)
        if raw is None:
            continue

        origin = getattr(f.type, "__origin__", None)
        if origin is list:
            item_type = f.type.__args__[0]
            if hasattr(item_type, "__dataclass_fields__"):
                kwargs[f.name] = [_parse(item_type, item) for item in raw]
            else:
                kwargs[f.name] = raw
        elif hasattr(f.type, "__dataclass_fields__"):
            kwargs[f.name] = _parse(f.type, raw)
        else:
            kwargs[f.name] = raw

    return cls(**kwargs)


# ── API client ───────────────────────────────────────────────────────────────


class Client:
    def __init__(self, token: str):
        self.token = token
        self.session = requests.Session()
        self.session.headers.update(
            {
                "Accept": "application/vnd.github+json",
                "Authorization": f"Bearer {token}",
                "X-GitHub-Api-Version": API_VERSION,
            }
        )

    def fetch_metrics(
        self, org: str, since: str = "", until: str = ""
    ) -> list[DayMetrics]:
        all_metrics: list[DayMetrics] = []
        page = 1

        while True:
            params: dict[str, str | int] = {"per_page": PER_PAGE, "page": page}
            if since:
                params["since"] = since
            if until:
                params["until"] = until

            resp = self.session.get(
                f"{BASE_URL}/orgs/{org}/copilot/metrics", params=params
            )

            if resp.status_code != 200:
                raise SystemExit(
                    f"Error: API returned {resp.status_code}: {resp.text}"
                )

            items = resp.json()
            for item in items:
                all_metrics.append(_parse(DayMetrics, item))

            if len(items) < PER_PAGE:
                break
            page += 1

        return all_metrics
