# 配置
TIMESTAMP_FORMAT := %Y%m%d-%H%M%S
TAG_FILE := .env
BACKUP_SUFFIX := .bak

# 命令
CURRENT_TIME := $(shell date +$(TIMESTAMP_FORMAT))
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null | tr '/' '-')
GIT_USERNAME := $(shell git config user.name 2>/dev/null || whoami)
CHECK_GIT := $(shell git rev-parse --is-inside-work-tree 2>/dev/null)

# 颜色输出
YELLOW := \033[1;33m
GREEN := \033[1;32m
RED := \033[1;31m
NC := \033[0m # No Color

.PHONY: tag
tag:
	@if [ -z "$(CHECK_GIT)" ]; then \
		echo "$(RED)Error: Not in a git repository$(NC)"; \
		exit 1; \
	fi
	@if [ -z "$(CURRENT_BRANCH)" ]; then \
		echo "$(RED)Error: Could not determine current branch$(NC)"; \
		exit 1; \
	fi
	@if [ -z "$(GIT_USERNAME)" ]; then \
		echo "$(RED)Error: Could not determine username$(NC)"; \
		exit 1; \
	fi
	@if [ -f "$(TAG_FILE)" ]; then \
		cp "$(TAG_FILE)" "$(TAG_FILE)$(BACKUP_SUFFIX)"; \
		echo "$(YELLOW)Backup created: $(TAG_FILE)$(BACKUP_SUFFIX)$(NC)"; \
	fi
	@echo "USER_ImageTag=$(GIT_USERNAME).$(CURRENT_BRANCH).$(CURRENT_TIME)" > "$(TAG_FILE)"
	@if [ $$? -eq 0 ]; then \
		echo "$(GREEN)Successfully generated tag in $(TAG_FILE):$(NC)"; \
		sed 's/^USER_ImageTag=//' "$(TAG_FILE)"; \
	else \
		echo "$(RED)Error: Failed to write to $(TAG_FILE)$(NC)"; \
		if [ -f "$(TAG_FILE)$(BACKUP_SUFFIX)" ]; then \
			mv "$(TAG_FILE)$(BACKUP_SUFFIX)" "$(TAG_FILE)"; \
			echo "$(YELLOW)Restored from backup$(NC)"; \
		fi; \
		exit 1; \
	fi
