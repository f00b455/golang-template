// Terminal News Reader - JavaScript Implementation
(function() {
    'use strict';

    // Configuration
    const CONFIG = {
        API_ENDPOINT: '/api/rss/spiegel/top5',
        EXPORT_ENDPOINT: '/api/rss/spiegel/export',
        REFRESH_INTERVAL: 300000, // 5 minutes
        FILTER_DELAY: 50, // 50ms debounce for real-time filtering
        THEMES: ['default', 'amber', 'matrix'],
        CACHE_KEY: 'terminal_rss_cache',
        CACHE_DURATION: 3600000 // 1 hour
    };

    // State management
    const state = {
        rssItems: [],
        filteredItems: [],
        selectedIndex: 0,
        commandHistory: [],
        historyIndex: -1,
        currentTheme: 'default',
        isOnline: navigator.onLine,
        filterTimer: null
    };

    // DOM Elements
    const elements = {
        commandInput: null,
        output: null,
        rssContainer: null,
        feedCount: null,
        datetime: null,
        position: null,
        filterStatus: null,
        suggestions: null,
        onlineStatus: null,
        exportJsonBtn: null,
        exportCsvBtn: null,
        exportStatus: null
    };

    // Initialize on DOM ready
    document.addEventListener('DOMContentLoaded', init);

    function init() {
        cacheElements();
        setupEventListeners();
        initializeTerminal();
        loadRSSFeed();
        startClock();
        setupMatrixRain();
        checkOnlineStatus();
    }

    function cacheElements() {
        elements.commandInput = document.getElementById('command-input');
        elements.output = document.getElementById('output');
        elements.rssContainer = document.getElementById('rss-container');
        elements.feedCount = document.getElementById('feed-count');
        elements.datetime = document.getElementById('datetime');
        elements.position = document.getElementById('position');
        elements.filterStatus = document.getElementById('filter-status');
        elements.suggestions = document.getElementById('suggestions');
        elements.onlineStatus = document.querySelector('.online-status');
        elements.exportJsonBtn = document.getElementById('export-json');
        elements.exportCsvBtn = document.getElementById('export-csv');
        elements.exportStatus = document.getElementById('export-status');
    }

    function setupEventListeners() {
        // Command input
        elements.commandInput.addEventListener('input', handleInput);
        elements.commandInput.addEventListener('keydown', handleKeydown);

        // Export buttons
        elements.exportJsonBtn.addEventListener('click', () => exportDataFormat('json'));
        elements.exportCsvBtn.addEventListener('click', () => exportDataFormat('csv'));

        // Global keyboard shortcuts
        document.addEventListener('keydown', handleGlobalKeys);

        // Online/offline detection
        window.addEventListener('online', () => updateOnlineStatus(true));
        window.addEventListener('offline', () => updateOnlineStatus(false));

        // Auto-refresh
        setInterval(loadRSSFeed, CONFIG.REFRESH_INTERVAL);
    }

    function initializeTerminal() {
        typewriterEffect('Initializing Terminal News Reader...', () => {
            setTimeout(() => {
                elements.output.style.display = 'none';
                elements.rssContainer.classList.remove('hidden');
                displaySystemMessage('System ready. Type :help for commands.');
            }, 500);
        });
    }

    // RSS Feed Management
    async function loadRSSFeed() {
        try {
            // Try to load from cache first if offline
            if (!state.isOnline) {
                const cached = loadFromCache();
                if (cached) {
                    state.rssItems = cached;
                    renderRSSItems();
                    displaySystemMessage('Loaded from cache (offline mode)');
                    return;
                }
            }

            const response = await fetch(CONFIG.API_ENDPOINT);
            if (!response.ok) throw new Error(`HTTP ${response.status}`);

            const data = await response.json();
            state.rssItems = processRSSData(data);
            state.filteredItems = [...state.rssItems];

            // Cache the data
            saveToCache(state.rssItems);

            renderRSSItems();
            updateFeedCount();
            displaySystemMessage(`Loaded ${state.rssItems.length} items`);

        } catch (error) {
            console.error('Failed to load RSS feed:', error);

            // Try cache as fallback
            const cached = loadFromCache();
            if (cached) {
                state.rssItems = cached;
                renderRSSItems();
                displaySystemMessage('Failed to fetch. Using cached data.', 'warning');
            } else {
                displaySystemMessage(`Error: ${error.message}`, 'error');
            }
        }
    }

    function processRSSData(data) {
        // Handle different API response formats
        if (data.headlines) {
            return data.headlines.map(item => ({
                title: item.title,
                link: item.link,
                description: item.description || '',
                publishedAt: item.publishedAt || item.published_at,
                source: item.source || 'Unknown'
            }));
        } else if (data.items) {
            return data.items;
        } else if (Array.isArray(data)) {
            return data;
        }
        return [];
    }

    function renderRSSItems() {
        const container = elements.rssContainer;
        container.innerHTML = '';

        const itemsToRender = state.filteredItems.length > 0 ?
            state.filteredItems : state.rssItems;

        itemsToRender.forEach((item, index) => {
            const article = createRSSElement(item, index);
            container.appendChild(article);

            // Stagger animation
            setTimeout(() => {
                article.style.opacity = '1';
            }, index * 50);
        });

        updatePosition();
    }

    function createRSSElement(item, index) {
        const article = document.createElement('article');
        article.className = 'rss-item';
        article.dataset.index = index;

        article.innerHTML = `
            <a href="${item.link}" target="_blank" class="rss-title">
                [${index + 1}] ${escapeHtml(item.title)}
            </a>
            <div class="rss-meta">
                <span class="date">${formatDate(item.publishedAt)}</span>
                <span class="source">${escapeHtml(item.source)}</span>
            </div>
            ${item.description ? `<div class="rss-description">${escapeHtml(item.description)}</div>` : ''}
        `;

        article.addEventListener('click', (e) => {
            if (!e.target.classList.contains('rss-title')) {
                e.preventDefault();
                selectItem(index);
            }
        });

        return article;
    }

    // Filtering System
    function handleInput(e) {
        const value = e.target.value;

        // Clear previous timer
        if (state.filterTimer) {
            clearTimeout(state.filterTimer);
        }

        // Check for commands
        if (value.startsWith(':')) {
            handleCommand(value);
            return;
        }

        // Debounced filtering for performance
        state.filterTimer = setTimeout(() => {
            applyFilter(value);
        }, CONFIG.FILTER_DELAY);
    }

    function applyFilter(query) {
        if (!query) {
            state.filteredItems = [...state.rssItems];
            elements.filterStatus.textContent = '';
        } else {
            const filters = parseFilterQuery(query);
            state.filteredItems = state.rssItems.filter(item =>
                matchesFilters(item, filters)
            );
            elements.filterStatus.textContent = `Filter: "${query}" (${state.filteredItems.length} matches)`;
        }

        renderRSSItems();
    }

    function parseFilterQuery(query) {
        const filters = {
            include: [],
            exclude: [],
            exact: [],
            regex: null
        };

        // Parse advanced filter syntax
        const tokens = query.match(/(\+\w+|-\w+|"[^"]+"|\/[^\/]+\/|\S+)/g) || [];

        tokens.forEach(token => {
            if (token.startsWith('+')) {
                filters.include.push(token.slice(1).toLowerCase());
            } else if (token.startsWith('-')) {
                filters.exclude.push(token.slice(1).toLowerCase());
            } else if (token.startsWith('"') && token.endsWith('"')) {
                filters.exact.push(token.slice(1, -1).toLowerCase());
            } else if (token.startsWith('/') && token.endsWith('/')) {
                try {
                    filters.regex = new RegExp(token.slice(1, -1), 'i');
                } catch (e) {
                    console.error('Invalid regex:', e);
                }
            } else {
                filters.include.push(token.toLowerCase());
            }
        });

        return filters;
    }

    function matchesFilters(item, filters) {
        const text = `${item.title} ${item.description} ${item.source}`.toLowerCase();

        // Check excludes first
        for (const exclude of filters.exclude) {
            if (text.includes(exclude)) return false;
        }

        // Check includes
        for (const include of filters.include) {
            if (!text.includes(include)) return false;
        }

        // Check exact matches
        for (const exact of filters.exact) {
            if (!text.includes(exact)) return false;
        }

        // Check regex
        if (filters.regex && !filters.regex.test(text)) {
            return false;
        }

        return true;
    }

    // Command System
    function handleCommand(command) {
        const cmd = command.toLowerCase().trim();
        const parts = cmd.split(' ');
        const mainCommand = parts[0];
        const args = parts.slice(1);

        switch (mainCommand) {
            case ':help':
                showHelp();
                break;
            case ':refresh':
                loadRSSFeed();
                elements.commandInput.value = '';
                break;
            case ':clear':
                clearScreen();
                break;
            case ':theme':
                cycleTheme();
                break;
            case ':stats':
                showStats();
                break;
            case ':export':
                // Support :export json or :export csv
                if (args[0] === 'json' || args[0] === 'csv') {
                    exportDataFormat(args[0]);
                    elements.commandInput.value = '';
                } else {
                    exportData();
                }
                break;
            case ':vim':
                enableVimMode();
                break;
            default:
                displaySystemMessage(`Unknown command: ${cmd}`, 'error');
        }
    }

    function showHelp() {
        const helpText = `
Available Commands:
  :help         - Show this help message
  :refresh      - Reload RSS feed
  :clear        - Clear the screen
  :theme        - Cycle through themes
  :stats        - Show statistics
  :export       - Export as JSON (default)
  :export json  - Export as JSON format
  :export csv   - Export as CSV format
  :vim          - Toggle vim keybindings

Filter Syntax:
  word      - Include items containing 'word'
  +word     - Must include 'word'
  -word     - Must not include 'word'
  "phrase"  - Exact phrase match
  /regex/   - Regular expression match

Keyboard Shortcuts:
  j/↓       - Next item
  k/↑       - Previous item
  /         - Focus search
  Enter     - Open selected item
  Escape    - Clear filter
  Tab       - Autocomplete`;

        displaySystemMessage(helpText);
        elements.commandInput.value = '';
    }

    function clearScreen() {
        elements.rssContainer.innerHTML = '';
        setTimeout(() => {
            renderRSSItems();
            displaySystemMessage('Screen cleared');
        }, 100);
        elements.commandInput.value = '';
    }

    function cycleTheme() {
        const currentIndex = CONFIG.THEMES.indexOf(state.currentTheme);
        const nextIndex = (currentIndex + 1) % CONFIG.THEMES.length;
        state.currentTheme = CONFIG.THEMES[nextIndex];

        document.body.className = state.currentTheme === 'default' ? '' : `theme-${state.currentTheme}`;
        displaySystemMessage(`Theme changed to: ${state.currentTheme}`);
        elements.commandInput.value = '';
    }

    function showStats() {
        const stats = `
Statistics:
  Total items: ${state.rssItems.length}
  Filtered items: ${state.filteredItems.length}
  Cache size: ${getCacheSize()} KB
  Online status: ${state.isOnline ? 'Connected' : 'Offline'}
  Theme: ${state.currentTheme}
  Last refresh: ${new Date().toLocaleTimeString()}`;

        displaySystemMessage(stats);
        elements.commandInput.value = '';
    }

    async function exportDataFormat(format) {
        try {
            // Show exporting status
            updateExportStatus('Exporting...', 'success');

            // Disable buttons during export
            elements.exportJsonBtn.disabled = true;
            elements.exportCsvBtn.disabled = true;

            // Get current filter from input
            const filterValue = elements.commandInput.value;
            const hasFilter = filterValue && !filterValue.startsWith(':');

            // Determine items to export
            const itemsToExport = hasFilter ? state.filteredItems : state.rssItems;

            // Build URL with query params
            const url = new URL(CONFIG.EXPORT_ENDPOINT, window.location.origin);
            url.searchParams.append('format', format);

            // If we have items to export, include limit
            if (itemsToExport.length > 0) {
                url.searchParams.append('limit', itemsToExport.length);
            }

            // Make API call
            const response = await fetch(url.toString());

            if (!response.ok) {
                throw new Error(`Export failed: ${response.status}`);
            }

            // Get filename from Content-Disposition header or generate one
            const contentDisposition = response.headers.get('Content-Disposition');
            let filename = `rss-export-${Date.now()}.${format}`;

            if (contentDisposition) {
                const filenameMatch = contentDisposition.match(/filename="?(.+?)"?(?:;|$)/);
                if (filenameMatch) {
                    filename = filenameMatch[1];
                }
            }

            // Add filter to filename if present
            if (hasFilter) {
                const sanitizedFilter = filterValue.replace(/[^a-z0-9]/gi, '_').substring(0, 20);
                filename = filename.replace(`.${format}`, `-${sanitizedFilter}.${format}`);
            }

            // Get blob from response
            const blob = await response.blob();

            // Create download link
            const downloadUrl = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = downloadUrl;
            a.download = filename;

            // Trigger download
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);

            // Clean up
            URL.revokeObjectURL(downloadUrl);

            // Show success message
            const message = hasFilter ?
                `Exported ${itemsToExport.length} filtered items as ${format.toUpperCase()}` :
                `Exported ${itemsToExport.length} items as ${format.toUpperCase()}`;

            updateExportStatus(message, 'success');
            displaySystemMessage(message);

        } catch (error) {
            console.error('Export failed:', error);
            updateExportStatus(`Export failed: ${error.message}`, 'error');
            displaySystemMessage(`Export failed: ${error.message}`, 'error');
        } finally {
            // Re-enable buttons
            elements.exportJsonBtn.disabled = false;
            elements.exportCsvBtn.disabled = false;

            // Clear status after 3 seconds
            setTimeout(() => {
                updateExportStatus('', '');
            }, 3000);
        }
    }

    function updateExportStatus(message, type) {
        if (elements.exportStatus) {
            elements.exportStatus.textContent = message;
            elements.exportStatus.className = 'export-status';
            if (type) {
                elements.exportStatus.classList.add(type);
            }
        }
    }

    function exportData() {
        // Legacy function for :export command
        exportDataFormat('json');
        elements.commandInput.value = '';
    }

    // Keyboard Navigation
    function handleKeydown(e) {
        switch (e.key) {
            case 'ArrowUp':
                e.preventDefault();
                navigateHistory(-1);
                break;
            case 'ArrowDown':
                e.preventDefault();
                navigateHistory(1);
                break;
            case 'Tab':
                e.preventDefault();
                autocomplete();
                break;
            case 'Enter':
                if (e.target.value.startsWith(':')) {
                    state.commandHistory.push(e.target.value);
                    state.historyIndex = state.commandHistory.length;
                }
                break;
            case 'Escape':
                e.preventDefault();
                clearFilter();
                break;
        }
    }

    function handleGlobalKeys(e) {
        // Skip if input is focused
        if (document.activeElement === elements.commandInput) return;

        // Handle Ctrl+E shortcuts for export
        if (e.ctrlKey && e.key === 'e') {
            e.preventDefault();
            // Wait for next key
            const handleExportKey = (evt) => {
                if (evt.key === 'j') {
                    exportDataFormat('json');
                } else if (evt.key === 'c') {
                    exportDataFormat('csv');
                }
                document.removeEventListener('keydown', handleExportKey);
            };
            document.addEventListener('keydown', handleExportKey);
            displaySystemMessage('Export mode: Press J for JSON, C for CSV');
            return;
        }

        switch (e.key) {
            case 'j':
            case 'ArrowDown':
                e.preventDefault();
                navigateItems(1);
                break;
            case 'k':
            case 'ArrowUp':
                e.preventDefault();
                navigateItems(-1);
                break;
            case '/':
                e.preventDefault();
                elements.commandInput.focus();
                break;
            case 'Enter':
                e.preventDefault();
                openSelectedItem();
                break;
        }
    }

    function navigateItems(direction) {
        const items = document.querySelectorAll('.rss-item');
        if (items.length === 0) return;

        // Remove previous selection
        items[state.selectedIndex]?.classList.remove('selected');

        // Update index
        state.selectedIndex = Math.max(0,
            Math.min(items.length - 1, state.selectedIndex + direction));

        // Add new selection
        items[state.selectedIndex]?.classList.add('selected');
        items[state.selectedIndex]?.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        });

        updatePosition();
    }

    function selectItem(index) {
        const items = document.querySelectorAll('.rss-item');
        items.forEach(item => item.classList.remove('selected'));

        state.selectedIndex = index;
        items[index]?.classList.add('selected');
        updatePosition();
    }

    function openSelectedItem() {
        const items = state.filteredItems.length > 0 ?
            state.filteredItems : state.rssItems;
        const item = items[state.selectedIndex];

        if (item && item.link) {
            window.open(item.link, '_blank');
            displaySystemMessage(`Opening: ${item.title}`);
        }
    }

    function clearFilter() {
        elements.commandInput.value = '';
        applyFilter('');
        displaySystemMessage('Filter cleared');
    }

    // History Navigation
    function navigateHistory(direction) {
        if (state.commandHistory.length === 0) return;

        state.historyIndex = Math.max(0,
            Math.min(state.commandHistory.length - 1, state.historyIndex + direction));

        elements.commandInput.value = state.commandHistory[state.historyIndex] || '';
    }

    // Autocomplete
    function autocomplete() {
        const value = elements.commandInput.value;
        if (!value) return;

        const suggestions = [];

        // Command suggestions
        if (value.startsWith(':')) {
            const commands = [':help', ':refresh', ':clear', ':theme', ':stats', ':export', ':vim'];
            suggestions.push(...commands.filter(cmd => cmd.startsWith(value)));
        } else {
            // Content suggestions from RSS items
            const words = new Set();
            state.rssItems.forEach(item => {
                const text = `${item.title} ${item.source}`.toLowerCase();
                text.split(/\s+/).forEach(word => {
                    if (word.length > 3 && word.startsWith(value.toLowerCase())) {
                        words.add(word);
                    }
                });
            });
            suggestions.push(...Array.from(words).slice(0, 5));
        }

        if (suggestions.length > 0) {
            elements.commandInput.value = suggestions[0];
        }
    }

    // Utility Functions
    function typewriterEffect(text, callback) {
        const output = elements.output.querySelector('.loading-animation');
        if (!output) return;

        let index = 0;
        const interval = setInterval(() => {
            if (index < text.length) {
                output.textContent = text.slice(0, index + 1) + '_';
                index++;
            } else {
                clearInterval(interval);
                if (callback) callback();
            }
        }, 50);
    }

    function displaySystemMessage(message, type = 'success') {
        const msgElement = document.createElement('div');
        msgElement.className = `system-message ${type}`;
        msgElement.textContent = message;

        if (elements.output.style.display === 'none') {
            elements.output.style.display = 'block';
            elements.output.innerHTML = '';
        }

        elements.output.appendChild(msgElement);
        elements.output.scrollTop = elements.output.scrollHeight;

        // Auto-hide after 5 seconds for non-error messages
        if (type !== 'error') {
            setTimeout(() => {
                msgElement.style.opacity = '0';
                setTimeout(() => msgElement.remove(), 500);
            }, 5000);
        }
    }

    function updateFeedCount() {
        elements.feedCount.textContent = state.rssItems.length;
    }

    function updatePosition() {
        const total = state.filteredItems.length || state.rssItems.length;
        elements.position.textContent = total > 0 ?
            `${state.selectedIndex + 1}/${total}` : '0/0';
    }

    function startClock() {
        const updateClock = () => {
            const now = new Date();
            elements.datetime.textContent = now.toLocaleString();
        };

        updateClock();
        setInterval(updateClock, 1000);
    }

    function updateOnlineStatus(isOnline) {
        state.isOnline = isOnline;
        elements.onlineStatus.textContent = isOnline ? 'ONLINE' : 'OFFLINE';
        elements.onlineStatus.className = isOnline ? 'online-status' : 'online-status error';

        if (isOnline) {
            displaySystemMessage('Connection restored', 'success');
            loadRSSFeed();
        } else {
            displaySystemMessage('Connection lost - working offline', 'warning');
        }
    }

    function checkOnlineStatus() {
        updateOnlineStatus(navigator.onLine);
    }

    // Cache Management
    function saveToCache(data) {
        try {
            const cacheData = {
                data: data,
                timestamp: Date.now()
            };
            localStorage.setItem(CONFIG.CACHE_KEY, JSON.stringify(cacheData));
        } catch (e) {
            console.error('Failed to save to cache:', e);
        }
    }

    function loadFromCache() {
        try {
            const cached = localStorage.getItem(CONFIG.CACHE_KEY);
            if (!cached) return null;

            const { data, timestamp } = JSON.parse(cached);
            const age = Date.now() - timestamp;

            if (age < CONFIG.CACHE_DURATION) {
                return data;
            }
        } catch (e) {
            console.error('Failed to load from cache:', e);
        }
        return null;
    }

    function getCacheSize() {
        const cache = localStorage.getItem(CONFIG.CACHE_KEY) || '';
        return Math.round(cache.length / 1024);
    }

    // Matrix Rain Effect
    function setupMatrixRain() {
        const canvas = document.getElementById('matrix-rain');
        if (!canvas) return;

        const ctx = canvas.getContext('2d');
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;

        const matrix = '01';
        const fontSize = 10;
        const columns = canvas.width / fontSize;
        const drops = Array(Math.floor(columns)).fill(1);

        function drawMatrix() {
            ctx.fillStyle = 'rgba(0, 0, 0, 0.04)';
            ctx.fillRect(0, 0, canvas.width, canvas.height);

            ctx.fillStyle = '#00ff00';
            ctx.font = fontSize + 'px monospace';

            for (let i = 0; i < drops.length; i++) {
                const text = matrix[Math.floor(Math.random() * matrix.length)];
                ctx.fillText(text, i * fontSize, drops[i] * fontSize);

                if (drops[i] * fontSize > canvas.height && Math.random() > 0.975) {
                    drops[i] = 0;
                }
                drops[i]++;
            }
        }

        setInterval(drawMatrix, 35);

        window.addEventListener('resize', () => {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        });
    }

    // Helper Functions
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    function formatDate(dateString) {
        if (!dateString) return 'Unknown';
        const date = new Date(dateString);
        return date.toLocaleString();
    }

    function enableVimMode() {
        displaySystemMessage('Vim mode enabled - use hjkl for navigation');
        document.addEventListener('keydown', function vimHandler(e) {
            if (document.activeElement === elements.commandInput) return;

            switch(e.key) {
                case 'h': // left
                case 'l': // right
                    // Could implement horizontal scrolling if needed
                    break;
                case 'g':
                    if (e.shiftKey) { // G - go to bottom
                        state.selectedIndex = (state.filteredItems.length || state.rssItems.length) - 1;
                        navigateItems(0);
                    } else { // gg - go to top
                        state.selectedIndex = 0;
                        navigateItems(0);
                    }
                    break;
            }
        });
        elements.commandInput.value = '';
    }

})();