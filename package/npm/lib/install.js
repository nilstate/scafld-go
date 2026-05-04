"use strict";

const fs = require("node:fs");
const https = require("node:https");
const path = require("node:path");
const { binaryPath, downloadURL, releaseVersion } = require("./platform");

if (process.env.SCAFLD_SKIP_DOWNLOAD === "1") {
  console.log("scafld: skipping native binary download because SCAFLD_SKIP_DOWNLOAD=1");
  process.exit(0);
}

if (releaseVersion() === "0.0.0" && !process.env.SCAFLD_INSTALL_BASE_URL) {
  console.log("scafld: development package version detected; skipping native binary download");
  process.exit(0);
}

const destination = binaryPath();
fs.mkdirSync(path.dirname(destination), { recursive: true });

download(downloadURL(), destination, 0).catch((err) => {
  console.error(`scafld: failed to install native binary: ${err.message}`);
  process.exit(1);
});

function download(url, destination, redirects) {
  if (redirects > 5) {
    return Promise.reject(new Error("too many redirects"));
  }

  return new Promise((resolve, reject) => {
    const tmp = `${destination}.tmp-${process.pid}`;
    const request = https.get(url, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        response.resume();
        resolve(download(response.headers.location, destination, redirects + 1));
        return;
      }

      if (response.statusCode !== 200) {
        response.resume();
        reject(new Error(`GET ${url} returned HTTP ${response.statusCode}`));
        return;
      }

      const file = fs.createWriteStream(tmp, { mode: 0o755 });
      response.pipe(file);
      file.on("finish", () => {
        file.close((err) => {
          if (err) {
            reject(err);
            return;
          }
          fs.chmodSync(tmp, 0o755);
          fs.renameSync(tmp, destination);
          console.log(`scafld: installed native binary ${destination}`);
          resolve();
        });
      });
      file.on("error", reject);
    });
    request.on("error", reject);
  });
}
