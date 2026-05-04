import os
import shutil
import stat
import tempfile
import urllib.request

from scafld_launcher import __version__
from scafld_launcher.platform import cache_root, target


def ensure_binary() -> str:
    if override := os.environ.get("SCAFLD_BINARY"):
        return override

    version = release_version()
    destination = os.path.join(cache_root(), version, f"scafld{target().ext}")
    if os.path.exists(destination):
        return destination

    if version == "0.0.0" and not os.environ.get("SCAFLD_INSTALL_BASE_URL"):
        raise RuntimeError("development package has no bundled scafld binary; set SCAFLD_BINARY")

    os.makedirs(os.path.dirname(destination), exist_ok=True)
    download(download_url(version), destination)
    return destination


def release_version() -> str:
    return os.environ.get("SCAFLD_INSTALL_VERSION", __version__).removeprefix("v")


def asset_name(version: str) -> str:
    selected = target()
    return f"scafld_{version.removeprefix('v')}_{selected.goos}_{selected.goarch}{selected.ext}"


def download_url(version: str) -> str:
    if base := os.environ.get("SCAFLD_INSTALL_BASE_URL"):
        return f"{base.rstrip('/')}/{asset_name(version)}"
    repo = os.environ.get("SCAFLD_GITHUB_REPOSITORY", "nilstate/scafld")
    return f"https://github.com/{repo}/releases/download/v{version}/{asset_name(version)}"


def download(url: str, destination: str) -> None:
    fd, tmp = tempfile.mkstemp(prefix=".scafld-", dir=os.path.dirname(destination))
    os.close(fd)
    try:
        with urllib.request.urlopen(url, timeout=60) as response, open(tmp, "wb") as out:
            shutil.copyfileobj(response, out)
        mode = os.stat(tmp).st_mode
        os.chmod(tmp, mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)
        os.replace(tmp, destination)
    except Exception:
        try:
            os.remove(tmp)
        except FileNotFoundError:
            pass
        raise
