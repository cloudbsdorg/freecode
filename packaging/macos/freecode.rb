class Freecode < Formula
  desc "Unified AI coding assistant"
  homepage "https://github.com/cloudbsdorg/freecode"
  url "https://github.com/cloudbsdorg/freecode/releases/download/v0.1.0/freecode-0.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "Unlimited"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags=-s -w", "-o", bin/"freecode", "./cmd/freecode"
  end

  test do
    system "#{bin}/freecode", "--version"
  end
end
