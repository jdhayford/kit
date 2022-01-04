# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kit < Formula
  desc ""
  homepage ""
  version "0.0.4"
  license "MIT"
  bottle :unneeded

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.4/kit_0.0.4_Darwin_x86_64.tar.gz"
      sha256 "cbbba2657f6f6490994fbba8cca132b7be67d1c8338d8e1faad246a25f498f08"
    end
    if Hardware::CPU.arm?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.4/kit_0.0.4_Darwin_arm64.tar.gz"
      sha256 "70154004e0216282ff5d7180cac4fbaf782739f8b6926dc9130354dbf4fcfab6"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.4/kit_0.0.4_Linux_x86_64.tar.gz"
      sha256 "4d5f39fdee418e67fd7be0ae4228325d5d5d93487b31df3d1727112b12979e51"
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/jdhayford/kit/releases/download/v0.0.4/kit_0.0.4_Linux_arm64.tar.gz"
      sha256 "33504b1900d29314cde14dc40a7f43251136e1399a222edc43a7e62b1f899b7a"
    end
  end

  depends_on "git"
  depends_on "zsh" => :optional

  def install
    bin.install "kit"
  end
end
