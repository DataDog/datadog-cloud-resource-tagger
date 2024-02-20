# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class DatadogCloudResourceTagger < Formula
  desc ""
  homepage ""
  version "0.0.5"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.5/datadog-cloud-resource-tagger_Darwin_arm64.tar.gz"
      sha256 "9b2bf25ecb985532c50a8af7ff76d515ad4c64264d37b81f9ad81f2940c71d5b"

      def install
        bin.install "datadog-cloud-resource-tagger"
        # Install shell completions
        generate_completions_from_executable(bin/"datadog-cloud-resource-tagger", "completion", shells: [:bash, :fish, :zsh], base_name: "datadog-cloud-resource-tagger")
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.5/datadog-cloud-resource-tagger_Darwin_x86_64.tar.gz"
      sha256 "2a8e9570dac77863f23eb6aca40575911e171e1671de712a6f7d9f56838cb15d"

      def install
        bin.install "datadog-cloud-resource-tagger"
        # Install shell completions
        generate_completions_from_executable(bin/"datadog-cloud-resource-tagger", "completion", shells: [:bash, :fish, :zsh], base_name: "datadog-cloud-resource-tagger")
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.5/datadog-cloud-resource-tagger_Linux_arm64.tar.gz"
      sha256 "14eef087e4efd2330194abade0a760d616002519ac1fc8ddf12a03d4182a1859"

      def install
        bin.install "datadog-cloud-resource-tagger"
        # Install shell completions
        generate_completions_from_executable(bin/"datadog-cloud-resource-tagger", "completion", shells: [:bash, :fish, :zsh], base_name: "datadog-cloud-resource-tagger")
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/DataDog/datadog-cloud-resource-tagger/releases/download/v0.0.5/datadog-cloud-resource-tagger_Linux_x86_64.tar.gz"
      sha256 "ed4f74a798d5f04710434d885ac39d3ba3c23f64811a2df57bae2cafbfb4707f"

      def install
        bin.install "datadog-cloud-resource-tagger"
        # Install shell completions
        generate_completions_from_executable(bin/"datadog-cloud-resource-tagger", "completion", shells: [:bash, :fish, :zsh], base_name: "datadog-cloud-resource-tagger")
      end
    end
  end
end
