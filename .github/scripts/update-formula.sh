#!/bin/bash
set -euo pipefail

VERSION="$1"
DARWIN_ARM64_URL="$2"
DARWIN_AMD64_URL="$3"
LINUX_AMD64_URL="$4"
LINUX_ARM64_URL="$5"

cat <<EOF
class Bleat < Formula
  desc "Find what's running on your ports"
  homepage "https://github.com/btbytes/bleat"
  version "$VERSION"

  on_macos do
    if Hardware::CPU.arm?
      url "$DARWIN_ARM64_URL"
      sha256 "$DARWIN_ARM64_SHA"
    else
      url "$DARWIN_AMD64_URL"
      sha256 "$DARWIN_AMD64_SHA"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "$LINUX_ARM64_URL"
      sha256 "$LINUX_ARM64_SHA"
    else
      url "$LINUX_AMD64_URL"
      sha256 "$LINUX_AMD64_SHA"
    end
  end

  def install
    bin.install "bleat"
  end

  test do
    system "#{bin}/bleat", "--help"
  end
end
EOF
