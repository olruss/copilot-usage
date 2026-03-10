import os
import subprocess


def get_token(flag_token: str | None) -> str:
    """Resolve a GitHub token.

    Priority: explicit flag > gh CLI > GITHUB_TOKEN env var.
    """
    if flag_token:
        return flag_token

    gh_token = _get_gh_token()
    if gh_token:
        return gh_token

    env_token = os.environ.get("GITHUB_TOKEN", "")
    if env_token:
        return env_token

    raise SystemExit(
        "Error: no GitHub token found\n\n"
        "To authenticate, do one of:\n"
        "  1. Install and login with gh CLI: gh auth login\n"
        "  2. Set GITHUB_TOKEN environment variable\n"
        "  3. Pass --token flag"
    )


def _get_gh_token() -> str | None:
    try:
        result = subprocess.run(
            ["gh", "auth", "token"],
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout.strip() or None
    except (subprocess.CalledProcessError, FileNotFoundError):
        return None
