#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');
const os = require('os');

// Package version
const VERSION = '0.1.0';
const REPO = 'kzeitar/freeport';

// Platform detection
const platform = os.platform();
const arch = os.arch();

// Map Node.js platforms to Go platforms
const platformMap = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

const archMap = {
  x64: 'amd64',
  arm64: 'arm64'
};

const goPlatform = platformMap[platform];
const goArch = archMap[arch];

if (!goPlatform || !goArch) {
  console.error(`Unsupported platform: ${platform}-${arch}`);
  process.exit(1);
}

// Binary name
const binaryName = platform === 'win32' ? 'freeport.exe' : 'freeport';
const releaseBinary = `freeport-${goPlatform}-${goArch}${platform === 'win32' ? '.exe' : ''}`;
const downloadUrl = `https://github.com/${REPO}/releases/download/v${VERSION}/${releaseBinary}`;

console.log(`Downloading freeport ${VERSION} for ${goPlatform}-${goArch}...`);
console.log(`URL: ${downloadUrl}`);

// Download the binary
const binDir = path.join(__dirname, 'bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

const outputPath = path.join(binDir, binaryName);

const file = fs.createWriteStream(outputPath);

https.get(downloadUrl, (response) => {
  if (response.statusCode === 302 || response.statusCode === 301) {
    https.get(response.headers.location, (redirectResponse) => {
      redirectResponse.pipe(file);
      file.on('finish', () => {
        file.close();
        fs.chmodSync(outputPath, 0o755);
        console.log('✓ Installed successfully');
      });
    }).on('error', (err) => {
      fs.unlink(outputPath, () => {});
      console.error('Error downloading binary:', err.message);
      process.exit(1);
    });
  } else {
    response.pipe(file);
    file.on('finish', () => {
      file.close();
      fs.chmodSync(outputPath, 0o755);
      console.log('✓ Installed successfully');
    });
  }
}).on('error', (err) => {
  fs.unlink(outputPath, () => {});
  console.error('Error downloading binary:', err.message);
  console.error('Please install via: go install github.com/kzeitar/freeport@latest');
  process.exit(1);
});
