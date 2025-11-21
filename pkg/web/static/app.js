// Inu Web UI - 主应用逻辑
(function () {
    'use strict';

    // 状态管理
    let credentials = null; // Basic Auth 凭据 (存内存)
    const STATE_KEY = 'inuState';

    // DOM 元素引用
    const elements = {
        // 脱敏视图
        anonymizeView: document.getElementById('anonymize-view'),
        entityTypes: document.getElementById('entity-types'),
        addCustomTypeBtn: document.getElementById('add-custom-type'),
        inputText: document.getElementById('input-text'),
        outputText: document.getElementById('output-text'),
        anonymizeBtn: document.getElementById('anonymize-btn'),
        switchToRestoreBtn: document.getElementById('switch-to-restore-btn'),

        // 还原视图
        restoreView: document.getElementById('restore-view'),
        entityMappingsDisplay: document.getElementById('entity-mappings-display'),
        anonymizedTextDisplay: document.getElementById('anonymized-text-display'),
        restoreInput: document.getElementById('restore-input'),
        restoreBtn: document.getElementById('restore-btn'),
        backToAnonymizeBtn: document.getElementById('back-to-anonymize-btn')
    };

    // ========== 初始化 ==========
    function init() {
        loadEntityTypesFromConfig();
        restoreStateFromSession();
        bindEvents();
    }

    // ========== 事件绑定 ==========
    function bindEvents() {
        elements.anonymizeBtn.addEventListener('click', handleAnonymize);
        elements.switchToRestoreBtn.addEventListener('click', switchToRestoreView);
        elements.restoreBtn.addEventListener('click', handleRestore);
        elements.backToAnonymizeBtn.addEventListener('click', switchToAnonymizeView);
        elements.addCustomTypeBtn.addEventListener('click', handleAddCustomType);
    }

    // ========== 配置加载 ==========
    async function loadEntityTypesFromConfig() {
        try {
            const response = await fetchWithAuth('/api/v1/config', {
                method: 'GET'
            });

            if (response.ok) {
                const config = await response.json();
                populateEntityTypes(config.entity_types || []);
            } else {
                // 如果配置端点不可用，使用默认类型
                console.warn('Config endpoint not available, using defaults');
                populateEntityTypes(['PERSON', 'ORG', 'EMAIL', 'PHONE', 'ADDRESS']);
            }
        } catch (error) {
            console.error('Failed to load config:', error);
            // 使用默认类型
            populateEntityTypes(['PERSON', 'ORG', 'EMAIL', 'PHONE', 'ADDRESS']);
        }
    }

    function populateEntityTypes(types) {
        elements.entityTypes.innerHTML = '';
        types.forEach(type => {
            const option = document.createElement('option');
            option.value = type;
            option.textContent = type;
            option.selected = true; // 默认选中所有类型
            elements.entityTypes.appendChild(option);
        });
    }

    // ========== 脱敏逻辑 ==========
    async function handleAnonymize() {
        const text = elements.inputText.value.trim();

        // 客户端验证
        if (!text) {
            showError(elements.outputText, '错误: 输入文本不能为空');
            elements.inputText.classList.add('error');
            return;
        }
        elements.inputText.classList.remove('error');

        // 获取选中的实体类型
        const selectedTypes = Array.from(elements.entityTypes.selectedOptions)
            .map(opt => opt.value);

        if (selectedTypes.length === 0) {
            showError(elements.outputText, '错误: 请至少选择一个实体类型');
            return;
        }

        // 显示加载状态
        setLoading(elements.anonymizeBtn, true);
        elements.outputText.textContent = '处理中...';
        elements.outputText.classList.add('loading');
        elements.outputText.classList.remove('error');

        try {
            const response = await fetchWithAuth('/api/v1/anonymize', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    text: text,
                    entity_types: selectedTypes
                })
            });

            if (response.ok) {
                const result = await response.json();

                // 显示结果
                elements.outputText.textContent = result.anonymized_text;
                elements.outputText.classList.remove('loading');

                // 保存状态到 sessionStorage
                saveStateToSession({
                    entities: result.entities,
                    anonymizedText: result.anonymized_text,
                    originalText: text,
                    entityTypes: selectedTypes
                });

                // 显示切换按钮
                elements.switchToRestoreBtn.style.display = 'inline-block';
            } else {
                const errorData = await response.json().catch(() => ({}));
                const errorMsg = errorData.error || `HTTP ${response.status}: ${response.statusText}`;
                showError(elements.outputText, `错误: ${errorMsg}`);
            }
        } catch (error) {
            showError(elements.outputText, '无法连接到服务器，请检查网络或服务器状态');
            console.error('Anonymize error:', error);
        } finally {
            setLoading(elements.anonymizeBtn, false);
            elements.outputText.classList.remove('loading');
        }
    }

    // ========== 还原逻辑 ==========
    async function handleRestore() {
        const text = elements.restoreInput.value.trim();

        // 客户端验证
        if (!text) {
            showError(elements.restoreInput.parentElement.querySelector('.output') || elements.restoreInput, '错误: 输入文本不能为空');
            elements.restoreInput.classList.add('error');
            return;
        }
        elements.restoreInput.classList.remove('error');

        // 从 sessionStorage 获取实体映射
        const state = loadStateFromSession();
        if (!state || !state.entities) {
            alert('错误: 未找到实体映射。请先执行脱敏操作。');
            return;
        }

        // 显示加载状态
        setLoading(elements.restoreBtn, true);
        const outputDiv = document.createElement('div');
        outputDiv.className = 'output loading';
        outputDiv.textContent = '处理中...';

        try {
            const response = await fetchWithAuth('/api/v1/restore', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    anonymized_text: text,
                    entities: state.entities
                })
            });

            if (response.ok) {
                const result = await response.json();
                elements.restoreInput.value = result.restored_text;
                elements.restoreInput.classList.remove('error');
            } else {
                const errorData = await response.json().catch(() => ({}));
                const errorMsg = errorData.error || `HTTP ${response.status}: ${response.statusText}`;
                alert(`还原失败: ${errorMsg}`);
            }
        } catch (error) {
            alert('无法连接到服务器，请检查网络或服务器状态');
            console.error('Restore error:', error);
        } finally {
            setLoading(elements.restoreBtn, false);
        }
    }

    // ========== 视图切换 ==========
    function switchToRestoreView() {
        const state = loadStateFromSession();
        if (!state) {
            alert('错误: 未找到脱敏数据');
            return;
        }

        // 显示实体映射
        displayEntityMappings(state.entities);

        // 显示脱敏文本
        elements.anonymizedTextDisplay.textContent = state.anonymizedText;

        // 切换视图
        elements.anonymizeView.style.display = 'none';
        elements.restoreView.style.display = 'block';
    }

    function switchToAnonymizeView() {
        const state = loadStateFromSession();

        // 恢复原始文本和脱敏结果
        if (state) {
            elements.inputText.value = state.originalText || '';
            elements.outputText.textContent = state.anonymizedText || '';
            elements.switchToRestoreBtn.style.display = 'inline-block';
        }

        // 切换视图
        elements.restoreView.style.display = 'none';
        elements.anonymizeView.style.display = 'block';
    }

    function displayEntityMappings(entities) {
        elements.entityMappingsDisplay.innerHTML = '';

        if (!entities || entities.length === 0) {
            elements.entityMappingsDisplay.textContent = '无实体映射';
            return;
        }

        entities.forEach(entity => {
            const item = document.createElement('div');
            item.className = 'mapping-item';
            // 显示 key -> values 格式
            const valuesStr = Array.isArray(entity.values) ? entity.values.join(', ') : String(entity.values);
            item.innerHTML = `
                <span class="placeholder">${escapeHtml(entity.key)}</span>
                <span class="arrow">→</span>
                <span class="value">${escapeHtml(valuesStr)}</span>
            `;
            elements.entityMappingsDisplay.appendChild(item);
        });
    }

    // ========== 自定义实体类型 ==========
    function handleAddCustomType() {
        const customType = prompt('输入自定义实体类型名称（例如: PRODUCT, LOCATION）:');
        if (customType && customType.trim()) {
            const type = customType.trim().toUpperCase();

            // 检查是否已存在
            const exists = Array.from(elements.entityTypes.options)
                .some(opt => opt.value === type);

            if (!exists) {
                const option = document.createElement('option');
                option.value = type;
                option.textContent = type;
                option.selected = true;
                elements.entityTypes.appendChild(option);
            }
        }
    }

    // ========== 状态管理 ==========
    function saveStateToSession(state) {
        try {
            sessionStorage.setItem(STATE_KEY, JSON.stringify(state));
        } catch (error) {
            console.error('Failed to save state:', error);
        }
    }

    function loadStateFromSession() {
        try {
            const data = sessionStorage.getItem(STATE_KEY);
            return data ? JSON.parse(data) : null;
        } catch (error) {
            console.error('Failed to load state:', error);
            return null;
        }
    }

    function restoreStateFromSession() {
        const state = loadStateFromSession();
        if (state) {
            // 恢复输入文本
            if (state.originalText) {
                elements.inputText.value = state.originalText;
            }

            // 恢复输出文本
            if (state.anonymizedText) {
                elements.outputText.textContent = state.anonymizedText;
                elements.switchToRestoreBtn.style.display = 'inline-block';
            }

            // 恢复实体类型选择
            if (state.entityTypes) {
                Array.from(elements.entityTypes.options).forEach(opt => {
                    opt.selected = state.entityTypes.includes(opt.value);
                });
            }
        }
    }

    // ========== HTTP 请求（带认证） ==========
    async function fetchWithAuth(url, options = {}) {
        // 如果有凭据，添加 Authorization 头
        if (credentials) {
            options.headers = options.headers || {};
            options.headers['Authorization'] = `Basic ${credentials}`;
        }

        let response = await fetch(url, options);

        // 处理 401 未授权
        if (response.status === 401 && !credentials) {
            const username = prompt('请输入用户名:');
            const password = prompt('请输入密码:');

            if (username && password) {
                credentials = btoa(`${username}:${password}`);
                options.headers = options.headers || {};
                options.headers['Authorization'] = `Basic ${credentials}`;
                response = await fetch(url, options);
            }
        }

        return response;
    }

    // ========== 工具函数 ==========
    function setLoading(button, isLoading) {
        if (isLoading) {
            button.disabled = true;
            button.dataset.originalText = button.textContent;
            button.innerHTML = '<span class="spinner"></span> 处理中...';
        } else {
            button.disabled = false;
            button.textContent = button.dataset.originalText || button.textContent;
        }
    }

    function showError(element, message) {
        element.textContent = message;
        element.classList.add('error');
        element.classList.remove('loading');
    }

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // ========== 启动应用 ==========
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
