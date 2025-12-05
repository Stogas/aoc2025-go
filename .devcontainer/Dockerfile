# any ARG values are default, override via build args in devcontainer.json instead
ARG GO_VERSION=1.25.4

FROM golang:${GO_VERSION} AS runner

ARG TARGETARCH TARGETOS

RUN apt-get update && apt-get install -y sudo bash git openssh-client curl nano jq gnupg2 tar time && rm -rf /var/lib/apt/lists/*

# installs to /usr/local/other-utils
# changest most rarely, so is installed first
RUN mkdir /usr/local/other-utils
# golangci-lint
ARG GOLANGCI_LINT_VERSION=2.6.2
RUN GOLANGCI_LINT_DOWNLOAD="golangci-lint-${GOLANGCI_LINT_VERSION}-${TARGETOS}-${TARGETARCH}" \
	&& curl -LO "https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/${GOLANGCI_LINT_DOWNLOAD}.tar.gz" \
	&& tar -xf "${GOLANGCI_LINT_DOWNLOAD}.tar.gz" -C /usr/local/other-utils --strip-components=1 "${GOLANGCI_LINT_DOWNLOAD}/golangci-lint"
# lefthook
RUN LATEST_LEFTHOOK=$(curl -s https://api.github.com/repos/evilmartians/lefthook/releases/latest | grep '"tag_name":' | cut -d'"' -f4 | sed 's/^v//') \
	&& case "$TARGETARCH" in \
	amd64) LEFTHOOK_ARCH="x86_64" ;; \
	arm64) LEFTHOOK_ARCH="aarch64" ;; \
	*) echo "Unsupported arch: $TARGETARCH"; exit 1 ;; \
	esac \
	&& LEFTHOOK_DOWNLOAD="lefthook_${LATEST_LEFTHOOK}_${TARGETOS}_${LEFTHOOK_ARCH}.gz" \
	&& curl -LO "https://github.com/evilmartians/lefthook/releases/download/v${LATEST_LEFTHOOK}/${LEFTHOOK_DOWNLOAD}" \
	&& gunzip -c "${LEFTHOOK_DOWNLOAD}" > /usr/local/other-utils/lefthook \
	&& chmod +x /usr/local/other-utils/lefthook

RUN useradd -m -s /bin/bash developer\
	&& echo "developer ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER developer
ENV PATH=$PATH:/home/developer/go/bin:/usr/local/go/bin:/usr/local/other-utils
