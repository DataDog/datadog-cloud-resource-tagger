# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class DatadogCloudResourceTagger < Formula
  desc ""
  homepage ""
  version "0.0.19"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.19/datadog-cloud-resource-tagger_Darwin_x86_64.tar.gz"
      sha256 "bfd6661523725bba27e6dbb718621ad3218ec45477233214541427b8cae72119"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.19/datadog-cloud-resource-tagger_Darwin_arm64.tar.gz"
      sha256 "008d25412757a9f6e3f61fa8687e7a48c3ecab9fd5b4ac651a39781f79750def"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.19/datadog-cloud-resource-tagger_Linux_arm64.tar.gz"
      sha256 "2bb41bdd96013b6cddcc64a2f8dcd71150174209044c0af48b1c251f9b03f620"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.19/datadog-cloud-resource-tagger_Linux_x86_64.tar.gz"
      sha256 "57fa23e87fea1fadd762264953a33dc6aa503b4eba5bb9b101bbada99555013d"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
  end
end
