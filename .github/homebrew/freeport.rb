# Homebrew formula for freeport
# Install via: brew install kzeitar/freeport/freeport
# Or from local file: brew install --formula freeport.rb

class Freeport < Formula
  desc "Command-line tool to free TCP ports by killing processes using them"
  homepage "https://github.com/kzeitar/freeport"
  url "https://github.com/kzeitar/freeport/archive/refs/tags/v0.1.0.tar.gz"
  sha256 :no_check # Will be computed automatically by Homebrew on first install
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"freeport"), "./cmd/freeport"
  end

  test do
    # Test that the binary runs
    assert_match(/freeport/, shell_output("#{bin}/freeport --help"))

    # Test that we can query an unused port
    port = free_port
    output = shell_output("#{bin}/freeport --list #{port}")
    assert_match(/No processes found/, output)
  end
end
