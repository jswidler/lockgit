require "rbconfig"
class Lockgit < Formula
  desc "A CLI tool for storing encrypted secrets in a git repo"
  homepage "https://github.com/jswidler/lockgit"
  version "0.5.0"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/jswidler/lockgit/releases/download/v0.5.0/lockgit_0.5.0_darwin_amd64.zip"
      sha256 "c9bffa00f6208b66e8c3549413951c0b57b89302c4ec66c0265a65c74664dc72"
    when /linux/
      url "https://github.com/jswidler/lockgit/releases/download/v0.5.0/lockgit_0.5.0_linux_amd64.tar.gz"
      sha256 "f5333f6a3bd70d42bf3a5ec48835b45c0ea3579b1b6ea775688d58f39cab83de"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  else
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/jswidler/lockgit/releases/download/v0.5.0/lockgit_0.5.0_darwin_386.zip"
      sha256 "a1446051985fd45a3f60d163529526b2fa4dcfda9428e7b75d4480641d15330d"
    when /linux/
      url "https://github.com/jswidler/lockgit/releases/download/v0.5.0/lockgit_0.5.0_linux_386.tar.gz"
      sha256 "2dc2eb2e05f9c133e776194b314fe79893f2ca68b29fa47ad1accc024d231830"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  end

  def install
    bin.install "lockgit"
  end

  test do
    system "#{bin}/lockgit"
  end
end
