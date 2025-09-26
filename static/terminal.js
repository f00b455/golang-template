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
        CACHE_DURATION: 3600000, // 1 hour
        ITEMS_PER_PAGE: 20, // Items to display per page
        MAX_ITEMS: 200, // Maximum items to fetch
        VIRTUAL_SCROLL_BUFFER: 5 // Items to render outside viewport
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
        filterTimer: null,
        currentPage: 1,
        totalPages: 1,
        isLoading: false,
        virtualScrollTop: 0,
        visibleRange: { start: 0, end: CONFIG.ITEMS_PER_PAGE },
        virtualScrollEnabled: false,
        itemHeight: 80 // Estimated height per item in pixels
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

        // Handle window resize for virtual scrolling
        window.addEventListener('resize', () => {
            if (state.virtualScrollEnabled) {
                renderRSSItems();
            }
        });
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
    async function loadRSSFeed(limit = CONFIG.MAX_ITEMS) {
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

            displayLoadingIndicator(true, 'Connecting to server...');
            const response = await fetch(`${CONFIG.API_ENDPOINT}?limit=${limit}`);
            if (!response.ok) throw new Error(`HTTP ${response.status}`);

            displayLoadingIndicator(true, 'Processing data...');
            const data = await response.json();
            state.rssItems = processRSSData(data);
            state.filteredItems = [...state.rssItems];

            displayLoadingIndicator(true, 'Caching data...');
            // Cache the data
            saveToCache(state.rssItems);

            displayLoadingIndicator(true, 'Rendering items...');
            updatePagination();
            renderRSSItems();
            updateFeedCount();
            displaySystemMessage(`Loaded ${state.rssItems.length} items`);
            displayLoadingIndicator(false);

        } catch (error) {
            console.error('Failed to load RSS feed:', error);

            // Try cache as fallback
            const cached = loadFromCache();
            if (cached) {
                state.rssItems = cached;
                updatePagination();
                renderRSSItems();
                displaySystemMessage('Failed to fetch. Using cached data.', 'warning');
                displayLoadingIndicator(false);
            } else {
                displaySystemMessage(`Error: ${error.message}`, 'error');
                displayLoadingIndicator(false);
            }
        }
    }

    // Pagination functions
    function updatePagination() {
        const items = state.filteredItems.length > 0 ? state.filteredItems : state.rssItems;
        state.totalPages = Math.ceil(items.length / CONFIG.ITEMS_PER_PAGE);

        // Ensure current page is valid
        if (state.currentPage > state.totalPages) {
            state.currentPage = Math.max(1, state.totalPages);
        }

        // Update pagination UI controls
        updatePaginationControls();
    }

    function updatePaginationStatus() {
        const items = state.filteredItems.length > 0 ? state.filteredItems : state.rssItems;
        const startItem = (state.currentPage - 1) * CONFIG.ITEMS_PER_PAGE + 1;
        const endItem = Math.min(state.currentPage * CONFIG.ITEMS_PER_PAGE, items.length);

        // Update status bar with pagination info in the format expected by BDD tests
        if (elements.position) {
            const statusText = items.length > CONFIG.ITEMS_PER_PAGE ?
                `${startItem}-${endItem} of ${items.length}` :
                `${items.length} items`;
            elements.position.textContent = statusText;
        }
    }

    // Add pagination controls UI
    function updatePaginationControls() {
        // Check if pagination controls container exists, if not create it
        let paginationContainer = document.getElementById('pagination-controls');
        if (!paginationContainer && elements.rssContainer) {
            paginationContainer = document.createElement('div');
            paginationContainer.id = 'pagination-controls';
            paginationContainer.className = 'pagination-controls';
            elements.rssContainer.parentNode.insertBefore(paginationContainer, elements.rssContainer.nextSibling);
        }

        if (!paginationContainer) return;

        const items = state.filteredItems.length > 0 ? state.filteredItems : state.rssItems;
        if (items.length <= CONFIG.ITEMS_PER_PAGE) {
            paginationContainer.style.display = 'none';
            return;
        }

        paginationContainer.style.display = 'flex';
        paginationContainer.innerHTML = `
            <button class="pagination-btn" id="page-first" ${state.currentPage === 1 ? 'disabled' : ''}>
                &laquo; First
            </button>
            <button class="pagination-btn" id="page-prev" ${state.currentPage === 1 ? 'disabled' : ''}>
                &lsaquo; Previous
            </button>
            <span class="pagination-info">
                Page <input type="number" id="page-input" min="1" max="${state.totalPages}"
                    value="${state.currentPage}" class="page-input"> of ${state.totalPages}
            </span>
            <button class="pagination-btn" id="page-next" ${state.currentPage === state.totalPages ? 'disabled' : ''}>
                Next &rsaquo;
            </button>
            <button class="pagination-btn" id="page-last" ${state.currentPage === state.totalPages ? 'disabled' : ''}>
                Last &raquo;
            </button>
        `;

        // Attach event handlers
        document.getElementById('page-first')?.addEventListener('click', () => navigateToPage(1));
        document.getElementById('page-prev')?.addEventListener('click', () => navigateToPage(state.currentPage - 1));
        document.getElementById('page-next')?.addEventListener('click', () => navigateToPage(state.currentPage + 1));
        document.getElementById('page-last')?.addEventListener('click', () => navigateToPage(state.totalPages));

        const pageInput = document.getElementById('page-input');
        if (pageInput) {
            pageInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    const page = parseInt(e.target.value);
                    if (page >= 1 && page <= state.totalPages) {
                        navigateToPage(page);
                    }
                }
            });
        }
    }

    function navigateToPage(pageNum) {
        if (pageNum < 1 || pageNum > state.totalPages) {
            displaySystemMessage(`Invalid page number. Valid range: 1-${state.totalPages}`, 'error');
            return;
        }

        state.currentPage = pageNum;
        renderRSSItems();
        displaySystemMessage(`Navigated to page ${pageNum}`);
    }

    // Add loading indicator functions
    function displayLoadingIndicator(show, progressText = null) {
        if (show) {
            state.isLoading = true;
            let loadingDiv = document.getElementById('loading-indicator');

            if (!loadingDiv) {
                loadingDiv = document.createElement('div');
                loadingDiv.id = 'loading-indicator';
                loadingDiv.className = 'loading-indicator';
                loadingDiv.innerHTML = `
                    <div class="loading-spinner"></div>
                    <div class="loading-text">Loading ${CONFIG.MAX_ITEMS} news items...</div>
                    <div class="loading-progress"></div>
                `;

                if (elements.rssContainer) {
                    elements.rssContainer.parentNode.insertBefore(loadingDiv, elements.rssContainer);
                }
            }

            // Update progress text if provided
            if (progressText && loadingDiv) {
                const progressElement = loadingDiv.querySelector('.loading-progress');
                if (progressElement) {
                    progressElement.textContent = progressText;
                }
            }
        } else {
            state.isLoading = false;
            const loadingDiv = document.getElementById('loading-indicator');
            if (loadingDiv) {
                loadingDiv.remove();
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
        const itemsToRender = state.filteredItems.length > 0 ?
            state.filteredItems : state.rssItems;

        // Enable virtual scrolling for large datasets
        if (itemsToRender.length > 100) {
            renderVirtualScrollItems(container, itemsToRender);
        } else {
            renderPaginatedItems(container, itemsToRender);
        }

        updatePosition();
        updatePaginationStatus();
    }

    function renderPaginatedItems(container, items) {
        container.innerHTML = '';

        // Calculate pagination
        const startIndex = (state.currentPage - 1) * CONFIG.ITEMS_PER_PAGE;
        const endIndex = Math.min(startIndex + CONFIG.ITEMS_PER_PAGE, items.length);
        const pageItems = items.slice(startIndex, endIndex);

        pageItems.forEach((item, index) => {
            const globalIndex = startIndex + index;
            const article = createRSSElement(item, globalIndex);
            container.appendChild(article);

            // Stagger animation
            setTimeout(() => {
                article.style.opacity = '1';
            }, index * 50);
        });
    }

    function renderVirtualScrollItems(container, items) {
        // Clear and setup virtual scroll container
        container.innerHTML = '';
        container.style.position = 'relative';
        container.style.overflowY = 'auto';
        container.style.height = '600px'; // Fixed height for virtual scroll

        // Create spacer for total height
        const totalHeight = items.length * state.itemHeight;
        const spacer = document.createElement('div');
        spacer.style.height = `${totalHeight}px`;
        spacer.style.position = 'relative';
        container.appendChild(spacer);

        // Create viewport for visible items
        const viewport = document.createElement('div');
        viewport.style.position = 'absolute';
        viewport.style.top = '0';
        viewport.style.left = '0';
        viewport.style.right = '0';
        spacer.appendChild(viewport);

        // Render initial visible items
        updateVirtualScroll(container, viewport, items);

        // Add scroll listener for virtual scrolling
        container.onscroll = () => {
            requestAnimationFrame(() => {
                updateVirtualScroll(container, viewport, items);
            });
        };

        state.virtualScrollEnabled = true;
    }

    function updateVirtualScroll(container, viewport, items) {
        const scrollTop = container.scrollTop;
        const containerHeight = container.clientHeight;

        // Calculate visible range with buffer
        const startIndex = Math.max(0, Math.floor(scrollTop / state.itemHeight) - CONFIG.VIRTUAL_SCROLL_BUFFER);
        const endIndex = Math.min(
            items.length,
            Math.ceil((scrollTop + containerHeight) / state.itemHeight) + CONFIG.VIRTUAL_SCROLL_BUFFER
        );

        // Only re-render if range changed significantly
        if (Math.abs(state.visibleRange.start - startIndex) > 1 ||
            Math.abs(state.visibleRange.end - endIndex) > 1) {

            state.visibleRange = { start: startIndex, end: endIndex };

            // Clear viewport and render visible items
            viewport.innerHTML = '';

            for (let i = startIndex; i < endIndex; i++) {
                const item = items[i];
                const article = createRSSElement(item, i);
                article.style.position = 'absolute';
                article.style.top = `${i * state.itemHeight}px`;
                article.style.left = '0';
                article.style.right = '0';
                viewport.appendChild(article);
            }
        }
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

        // Reset to first page when filtering
        state.currentPage = 1;
        updatePagination();
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
            case ':page':
                if (args.length > 0) {
                    const pageNum = parseInt(args[0]);
                    if (!isNaN(pageNum)) {
                        navigateToPage(pageNum);
                        elements.commandInput.value = '';
                    } else {
                        displaySystemMessage('Invalid page number. Usage: :page <number>', 'error');
                    }
                } else {
                    displaySystemMessage('Usage: :page <number> - Jump to specific page', 'info');
                }
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
            case ':first':
                navigateToPage(1);
                elements.commandInput.value = '';
                break;
            case ':last':
                navigateToPage(state.totalPages);
                elements.commandInput.value = '';
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
  :page <num>   - Jump to specific page
  :first        - Go to first page
  :last         - Go to last page

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
  Tab       - Autocomplete
  PageDown  - Next page
  PageUp    - Previous page
  Home      - First page
  End       - Last page
  1-9       - Jump to page (when not in input)`;

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

    // Constants for export functionality
    const EXPORT_STATUS_TIMEOUT = 3000; // ms
    const MAX_FILTER_LENGTH_DISPLAY = 20; // characters for filename
    const MAX_FILTER_LENGTH = 100; // max filter length
    let exportStatusTimeout = null;

    // Pure function for filename generation
    function generateExportFilename(format, filter, timestamp) {
        const base = `rss-export-${timestamp}.${format}`;
        if (!filter) return base;

        const sanitized = filter.replace(/[^a-z0-9]/gi, '_').substring(0, MAX_FILTER_LENGTH_DISPLAY);
        return base.replace(`.${format}`, `-${sanitized}.${format}`);
    }

    // Build export URL with parameters
    function buildExportUrl(format, limit) {
        const url = new URL(CONFIG.EXPORT_ENDPOINT, window.location.origin);
        url.searchParams.append('format', format);

        // Validate limit parameter
        if (limit > 0 && limit <= CONFIG.MAX_ITEMS) {
            url.searchParams.append('limit', limit);
        }

        return url;
    }

    // Extract filename from Content-Disposition header
    function extractFilename(contentDisposition, defaultFilename) {
        if (!contentDisposition) return defaultFilename;

        // More defensive parsing of Content-Disposition
        try {
            const filenameMatch = contentDisposition.match(/filename="?([^"\n\r;]+)"?/);
            return filenameMatch ? filenameMatch[1] : defaultFilename;
        } catch (e) {
            return defaultFilename;
        }
    }

    // Trigger file download
    function triggerDownload(blob, filename) {
        const downloadUrl = URL.createObjectURL(blob);

        try {
            const a = document.createElement('a');
            a.href = downloadUrl;
            a.download = filename;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } finally {
            // Ensure cleanup even if download fails
            URL.revokeObjectURL(downloadUrl);
        }
    }

    // Clear export status with proper timeout management
    function clearExportStatusDelayed() {
        // Cancel any existing timeout
        if (exportStatusTimeout) {
            clearTimeout(exportStatusTimeout);
        }

        exportStatusTimeout = setTimeout(() => {
            updateExportStatus('', '');
            exportStatusTimeout = null;
        }, EXPORT_STATUS_TIMEOUT);
    }

    // Validate export parameters
    function validateExportParams(format, filterValue) {
        if (!['json', 'csv'].includes(format)) {
            throw new Error('Invalid export format');
        }

        if (filterValue && filterValue.length > MAX_FILTER_LENGTH) {
            throw new Error(`Filter too long (max ${MAX_FILTER_LENGTH} characters)`);
        }
    }

    // Process export response
    async function processExportResponse(response, format, itemCount, hasFilter, filterValue) {
        if (!response.ok) {
            const errorMsg = response.status === 503 ? 'Service temporarily unavailable' :
                            response.status === 400 ? 'Invalid export parameters' :
                            `Export failed (${response.status})`;
            throw new Error(errorMsg);
        }

        const blob = await response.blob();
        const contentDisposition = response.headers.get('Content-Disposition');
        const defaultFilename = generateExportFilename(format, hasFilter ? filterValue : null, Date.now());
        const filename = extractFilename(contentDisposition, defaultFilename);

        triggerDownload(blob, filename);

        const message = `Exported ${itemCount}${hasFilter ? ' filtered' : ''} items as ${format.toUpperCase()}`;
        return message;
    }

    // Main export function (simplified)
    async function exportDataFormat(format) {
        try {
            // Get filter value
            const filterValue = elements.commandInput.value;
            const hasFilter = filterValue && !filterValue.startsWith(':');

            // Validate parameters
            validateExportParams(format, hasFilter ? filterValue : null);

            // Update UI
            updateExportStatus('Exporting...', 'success');
            elements.exportJsonBtn.disabled = true;
            elements.exportCsvBtn.disabled = true;

            // Build request
            const itemsToExport = hasFilter ? state.filteredItems : state.rssItems;
            const url = buildExportUrl(format, itemsToExport.length);

            // Make API call
            const response = await fetch(url.toString());

            // Process response
            const message = await processExportResponse(response, format, itemsToExport.length, hasFilter, filterValue);

            // Show success
            updateExportStatus(message, 'success');
            displaySystemMessage(message);

        } catch (error) {
            console.error('Export failed:', error);
            const errorMsg = error.message || 'Unknown error occurred';
            updateExportStatus(`Export failed: ${errorMsg}`, 'error');
            displaySystemMessage(`Export failed: ${errorMsg}`, 'error');
        } finally {
            elements.exportJsonBtn.disabled = false;
            elements.exportCsvBtn.disabled = false;
            clearExportStatusDelayed();
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

    // Track export key listener to prevent stacking
    let exportKeyListener = null;

    // Handle export mode
    function handleExportMode() {
        // Clean up any existing listener
        if (exportKeyListener) {
            document.removeEventListener('keydown', exportKeyListener);
            clearTimeout(exportKeyListener.timeoutId);
        }

        // Create new listener with timeout
        exportKeyListener = (evt) => {
            if (evt.key === 'j' || evt.key === 'J') {
                exportDataFormat('json');
                displaySystemMessage('Exporting as JSON...');
            } else if (evt.key === 'c' || evt.key === 'C') {
                exportDataFormat('csv');
                displaySystemMessage('Exporting as CSV...');
            } else {
                displaySystemMessage('Export cancelled');
            }

            document.removeEventListener('keydown', exportKeyListener);
            clearTimeout(exportKeyListener.timeoutId);
            exportKeyListener = null;
        };

        // Auto-cancel after 3 seconds
        exportKeyListener.timeoutId = setTimeout(() => {
            if (exportKeyListener) {
                document.removeEventListener('keydown', exportKeyListener);
                exportKeyListener = null;
                displaySystemMessage('Export mode cancelled (timeout)');
            }
        }, 3000);

        document.addEventListener('keydown', exportKeyListener);
        displaySystemMessage('Export mode: Press J for JSON, C for CSV');
    }

    function handleGlobalKeys(e) {
        // Skip if input is focused
        if (document.activeElement === elements.commandInput) return;

        // Handle Ctrl+E shortcuts for export
        if (e.ctrlKey && e.key === 'e') {
            e.preventDefault();
            handleExportMode();
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
            case 'PageDown':
                e.preventDefault();
                if (state.currentPage < state.totalPages) {
                    navigateToPage(state.currentPage + 1);
                }
                break;
            case 'PageUp':
                e.preventDefault();
                if (state.currentPage > 1) {
                    navigateToPage(state.currentPage - 1);
                }
                break;
            case 'Home':
                e.preventDefault();
                navigateToPage(1);
                break;
            case 'End':
                e.preventDefault();
                navigateToPage(state.totalPages);
                break;
            case '1':
            case '2':
            case '3':
            case '4':
            case '5':
            case '6':
            case '7':
            case '8':
            case '9':
                // Quick page navigation with number keys
                e.preventDefault();
                const pageNum = parseInt(e.key);
                if (pageNum <= state.totalPages) {
                    navigateToPage(pageNum);
                }
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
            const commands = [':help', ':refresh', ':clear', ':theme', ':stats', ':export', ':vim', ':page', ':first', ':last'];
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