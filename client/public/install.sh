#!/bin/bash
set -e

# UI Constants matching the CLI theme
RESET="\033[0m"
BOLD="\033[1m"
DIM="\033[2m"
CYAN="\033[36m"
GREEN="\033[32m"
RED="\033[31m"
YELLOW="\033[33m"

GREY="\033[37m"

print_banner() {
    echo ""
    echo -e "${GREY}  ██████╗  █████╗  ██████╗██╗  ██╗ █████╗  ██████╗ ███████╗${RESET}"
    echo -e "${GREY}  ██╔══██╗██╔══██╗██╔════╝██║ ██╔╝██╔══██╗██╔════╝ ██╔════╝${RESET}"
    echo -e "${GREY}  ██████╔╝███████║██║     █████╔╝ ███████║██║  ███╗█████╗  ${RESET}"
    echo -e "${GREY}  ██╔═══╝ ██╔══██║██║     ██╔═██╗ ██╔══██║██║   ██║██╔══╝  ${RESET}"
    echo -e "${GREY}  ██║     ██║  ██║╚██████╗██║  ██╗██║  ██║╚██████╔╝███████╗${RESET}"
    echo -e "${GREY}  ╚═╝     ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝${RESET}"
    echo ""
    echo -e "${GREY}  ███╗   ███╗ █████╗ ████████╗███████╗${RESET}"
    echo -e "${GREY}  ████╗ ████║██╔══██╗╚══██╔══╝██╔════╝${RESET}"
    echo -e "${GREY}  ██╔████╔██║███████║   ██║   █████╗  ${RESET}"
    echo -e "${GREY}  ██║╚██╔╝██║██╔══██║   ██║   ██╔══╝  ${RESET}"
    echo -e "${GREY}  ██║ ╚═╝ ██║██║  ██║   ██║   ███████╗${RESET}"
    echo -e "${GREY}  ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝${RESET}"
    echo ""
}

print_banner

OS="$(uname -s)"
ARCH="$(uname -m)"

# 1. Detect macOS
if [ "$OS" != "Darwin" ]; then
    echo -e "  ${RED}Error:${RESET} Package Mate Is Not Available on ${OS}"
    echo -e "         Only pristine macOS environments are currently supported."
    echo ""
    exit 1
fi

# 2. Detect Apple Silicon (M-Series)
if [ "$ARCH" != "arm64" ] && [ "$ARCH" != "aarch64" ]; then
    echo -e "  ${RED}Error:${RESET} Package Mate is exclusively optimized for Apple Silicon (M-Series)."
    echo -e "         Intel-based Macs are not supported natively."
    echo ""
    exit 1
fi

# 3. Check if Already Installed
if [ -f "/usr/local/bin/mate" ]; then
    echo -e "  ${YELLOW}Notice:${RESET}"
    echo ""
    echo -e "  ❯ Package Mate is already securely installed on your system."
    echo ""
    echo -e "  ❯ Run ${CYAN}mate${RESET} via the CLI to start managing your environment."
    echo ""
    exit 0
fi

# 4. Install Manager Engine (Homebrew) if missing
if ! command -v brew >/dev/null 2>&1; then
    echo -e "  ${DIM}Installing the Manager Engine (Homebrew)...${RESET}"
    echo ""
    # NONINTERACTIVE ensures it doesn't prompt for ENTER to continue installing
    NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" >/dev/null

    # Adding brew to path for future commands in this session
    if [ -d "/opt/homebrew/bin" ]; then
        export PATH="/opt/homebrew/bin:$PATH"
    fi
fi

# 5. Download and Extract the Binary
echo -e "  ${DIM}Fetching Package Mate v1.0.0...${RESET}"
echo ""
TARBALL_URL="https://github.com/yousefbustamiii/package-mate/releases/download/v1.0.0/mate-v1.0.0-darwin-arm64.tar.gz"

# Try to download the release binary
if ! curl -fsSL -o /tmp/mate.tar.gz "$TARBALL_URL" 2>/dev/null; then
    echo ""
    echo -e "  ${RED}Error:${RESET} Failed to download the release archive."
    echo -e "         The v1.0.0 release may not be available yet."
    echo ""
    echo -e "  ${YELLOW}Alternative:${RESET} Install from source instead:"
    echo ""
    echo -e "    ${DIM}git clone https://github.com/yousefbustamiii/package-mate.git${RESET}"
    echo -e "    ${DIM}cd package-mate${RESET}"
    echo -e "    ${DIM}go build -o mate .${RESET}"
    echo -e "    ${DIM}sudo mv mate /usr/local/bin/${RESET}"
    echo ""
    exit 1
fi

# Extract the tarball to /tmp (may contain mate binary in a subdirectory)
tar -xzf /tmp/mate.tar.gz -C /tmp/ || {
    echo ""
    echo -e "  ${RED}Error:${RESET} Failed to extract the mate binary from the downloaded archive."
    echo -e "         The release archive may be corrupted or in the wrong format."
    echo ""
    exit 1
}

# Find the mate binary (could be directly in /tmp or in a subdirectory)
if [ -f /tmp/mate ]; then
    MATE_BINARY="/tmp/mate"
elif [ -f /tmp/mate-v1.0.0-darwin-arm64/mate ]; then
    MATE_BINARY="/tmp/mate-v1.0.0-darwin-arm64/mate"
else
    # Try to find any 'mate' binary in the extracted contents
    MATE_BINARY=$(find /tmp -name "mate" -type f -perm /111 2>/dev/null | head -1)
fi

if [ -z "$MATE_BINARY" ] || [ ! -f "$MATE_BINARY" ]; then
    echo ""
    echo -e "  ${RED}Error:${RESET} Could not locate the mate binary in the archive."
    echo -e "         The release structure may have changed."
    echo ""
    exit 1
fi

# Ensure the binary is executable
chmod +x "$MATE_BINARY"

# Cleanup the tarball and extracted directory
rm -f /tmp/mate.tar.gz
rm -rf /tmp/mate-v1.0.0-darwin-arm64

# 6. Request sudo with standard elegant macOS style
echo -e "  ${GREEN}Download and extraction successful.${RESET}\n"
echo -e "  To finalize the installation, we need to move the binary to /usr/local/bin."
echo -e "  Please secure this action by verifying your identity."
echo ""
echo ""

# The native macOS terminal will automatically attach the key icon to the password prompt
sudo -p "  Password: " mv "$MATE_BINARY" /usr/local/bin/mate

echo ""
echo -e "  ${GREEN}Installation Complete.${RESET}"
echo -e "  Run ${CYAN}mate${RESET} to launch your developer experience."
echo ""
