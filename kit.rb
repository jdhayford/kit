# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kit < Formula
  desc ""
  homepage ""
  version "0.0.6"
  license "MIT"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.5/kit_0.0.5_Darwin_x86_64.tar.gz"
      sha256 "b797775190102bc39cb9610991869e79c1259ec93ac57a891773081748bebab0"
    end
    if Hardware::CPU.arm?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.5/kit_0.0.5_Darwin_arm64.tar.gz"
      sha256 "d36f89106fbbe60f5dfad16ee47c71b0ac647d67bc7b4ac17cc0ae7993eca1b0"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.5/kit_0.0.5_Linux_x86_64.tar.gz"
      sha256 "a2b6e36f334cacf3a256afdee0ce8eb8ce43c844081348a4233d0010ab49550b"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.5/kit_0.0.5_Linux_arm64.tar.gz"
      sha256 "5264edb6dff304e38a2f6c7b712a5092b42df8db6d21794abce42087732504a1"
    end
  end

  depends_on "git"
  depends_on "zsh" => :optional

  def install
    bin.install "kit"
  end
end
