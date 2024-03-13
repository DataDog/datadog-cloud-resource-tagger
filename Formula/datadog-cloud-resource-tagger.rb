# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class DatadogCloudResourceTagger < Formula
  desc ""
  homepage ""
  version "0.0.14"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.14/datadog-cloud-resource-tagger_Darwin_x86_64.tar.gz"
      sha256 "bc8371bcf05eedb1c60e0b49687aa952783f64e360e745315e247ab36f599380"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.14/datadog-cloud-resource-tagger_Darwin_arm64.tar.gz"
      sha256 "f28fc1aac56b9efde14e9f23b97a5f739843ec74f7d72e494576e2c63243f070"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.14/datadog-cloud-resource-tagger_Linux_arm64.tar.gz"
      sha256 "ca803b2cea84cab436e33fe160f47a25be396244653488d6d7a626437e9647e1"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.14/datadog-cloud-resource-tagger_Linux_x86_64.tar.gz"
      sha256 "76b353101a2423148205790c023144b72a5ccd492d5b72e3eb6d2af0f76cbcba"

      def install
        bin.install "datadog-cloud-resource-tagger"
      end
    end
  end
end
