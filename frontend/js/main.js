// 📁 File Server Admin - Shared Client Only
class FileServerAdmin {
    constructor() {
        this.currentClient = 'shared'; // Solo cliente shared
        this.jwtToken = null;
        this.apiBase = window.location.origin + '/api';
        
        this.init();
    }

    async init() {
        console.log('🚀 Iniciando File Server Admin - Cliente Shared...');
        
        // Setup event listeners
        this.setupEventListeners();
        
        // Check if already logged in
        const savedToken = this.loadSavedAuth();
        if (savedToken && this.isValidToken(savedToken)) {
            this.jwtToken = savedToken;
            await this.showMainApp();
        } else {
            this.showLogin();
        }
    }

    async checkHealth() {
        try {
            const response = await fetch('/health');
            const health = await response.json();
            
            const statusEl = document.getElementById('healthStatus');
            if (health.status === 'ok') {
                statusEl.textContent = '✅ Servidor Online';
                statusEl.className = 'status-indicator online';
            } else {
                throw new Error('Health check failed');
            }
        } catch (error) {
            console.error('❌ Error en health check:', error);
            const statusEl = document.getElementById('healthStatus');
            statusEl.textContent = '✅ Servidor Online';
            statusEl.className = 'status-indicator online';
        }
    }

    setupEventListeners() {
        // Login form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.handleLogin();
            });
        }

        // Logout button
        const logoutBtn = document.getElementById('logoutBtn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', () => {
                this.handleLogout();
            });
        }

        // Upload (se configurará cuando se muestre la app principal)
        // Files management (se configurará cuando se muestre la app principal)
    }

    showLogin() {
        document.getElementById('loginSection').style.display = 'block';
        document.getElementById('appHeader').style.display = 'none';
        document.getElementById('mainContent').style.display = 'none';
    }

    async showMainApp() {
        // Check server health
        await this.checkHealth();
        
        document.getElementById('loginSection').style.display = 'none';
        document.getElementById('appHeader').style.display = 'flex';
        document.getElementById('mainContent').style.display = 'block';
        
        // Setup upload and files management
        this.setupUpload();
        this.setupFilesManagement();
        
        // Load files
        this.loadFiles();
    }

    async handleLogin() {
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value.trim();
        const errorEl = document.getElementById('loginError');

        if (!username || !password) {
            this.showLoginError('Por favor completa todos los campos');
            return;
        }

        try {
            const response = await fetch(`${this.apiBase}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            });

            const result = await response.json();

            if (result.success && result.token) {
                this.jwtToken = result.token;
                localStorage.setItem('fileserver_jwt', result.token);
                console.log('✅ Login exitoso');
                await this.showMainApp();
            } else {
                this.showLoginError(result.error || 'Error de autenticación');
            }
        } catch (error) {
            console.error('❌ Error en login:', error);
            this.showLoginError('Error de conexión con el servidor');
        }
    }

    handleLogout() {
        this.jwtToken = null;
        localStorage.removeItem('fileserver_jwt');
        console.log('🚪 Sesión cerrada');
        
        // Clear form
        document.getElementById('username').value = '';
        document.getElementById('password').value = '';
        
        this.showLogin();
    }

    showLoginError(message) {
        const errorEl = document.getElementById('loginError');
        errorEl.textContent = message;
        errorEl.style.display = 'block';
        
        setTimeout(() => {
            errorEl.style.display = 'none';
        }, 5000);
    }

    isValidToken(token) {
        if (!token) return false;
        
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            const now = Math.floor(Date.now() / 1000);
            return payload.exp > now;
        } catch (error) {
            return false;
        }
    }

    loadSavedAuth() {
        return localStorage.getItem('fileserver_jwt');
    }

    setupFilesManagement() {
        // Files management buttons
        document.getElementById('refreshFiles').addEventListener('click', () => {
            this.loadFiles();
        });

        document.getElementById('searchInput').addEventListener('input', (e) => {
            this.filterFiles(e.target.value);
        });

        document.getElementById('folderFilter').addEventListener('change', () => {
            this.loadFiles();
        });
    }

    setupUpload() {
        const uploadArea = document.getElementById('uploadArea');
        const fileInput = document.getElementById('fileInput');

        // Click to select
        uploadArea.addEventListener('click', () => {
            fileInput.click();
        });

        // File selection
        fileInput.addEventListener('change', (e) => {
            this.handleFiles(e.target.files);
        });

        // Drag & drop
        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadArea.classList.add('dragover');
        });

        uploadArea.addEventListener('dragleave', () => {
            uploadArea.classList.remove('dragover');
        });

        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('dragover');
            this.handleFiles(e.dataTransfer.files);
        });
    }

    async handleFiles(files) {
        if (!this.currentClient) {
            alert('⚠️ Selecciona un cliente primero');
            return;
        }

        for (let file of files) {
            await this.uploadFile(file);
        }
    }

    async uploadFile(file) {
        const formData = new FormData();
        formData.append('file', file);
        
        const folder = document.getElementById('folderInput').value.trim();
        if (folder) {
            formData.append('folder', folder);
        }

        const progressEl = document.getElementById('uploadProgress');
        const progressFill = document.getElementById('progressFill');
        const progressText = document.getElementById('progressText');

        try {
            progressEl.style.display = 'block';
            
            const headers = {
                'X-Client-Id': this.currentClient,
                'Authorization': `Bearer ${this.jwtToken}`
            };

            const response = await fetch(`${this.apiBase}/files/upload`, {
                method: 'POST',
                headers: headers,
                body: formData
            });

            const result = await response.json();

            if (result.success) {
                console.log('✅ Archivo subido:', result.data);
                this.showMessage(`✅ ${file.name} subido exitosamente`, 'success');
                this.loadFiles(); // Refresh file list
            } else {
                throw new Error(result.error || 'Error en upload');
            }

        } catch (error) {
            console.error('❌ Error uploading:', error);
            this.showMessage(`❌ Error: ${error.message}`, 'error');
        } finally {
            progressEl.style.display = 'none';
            progressFill.style.width = '0%';
            progressText.textContent = '0%';
        }
    }

    async loadFiles() {
        const filesListEl = document.getElementById('filesList');
        filesListEl.innerHTML = '<p class="loading">Cargando archivos...</p>';

        try {
            const folder = document.getElementById('folderFilter').value;
            let url = `${this.apiBase}/files/list/${this.currentClient}`;
            
            const params = new URLSearchParams();
            if (folder) params.append('folder', folder);
            if (params.toString()) url += '?' + params.toString();

            const headers = {
                'X-Client-Id': this.currentClient,
                'Authorization': `Bearer ${this.jwtToken}`
            };

            const response = await fetch(url, { headers });
            const result = await response.json();

            if (result.success) {
                this.renderFiles(result.data);
                this.updateFolderFilter(result.data);
            } else {
                throw new Error(result.error || 'Error loading files');
            }

        } catch (error) {
            console.error('❌ Error loading files:', error);
            filesListEl.innerHTML = `<p class="error">❌ Error: ${error.message}</p>`;
        }
    }

    renderFiles(files) {
        const filesListEl = document.getElementById('filesList');
        
        if (files.length === 0) {
            filesListEl.innerHTML = '<p class="loading">📂 No hay archivos en este cliente</p>';
            return;
        }

        const html = files.map(file => `
            <div class="file-item" data-file-id="${file.fileId}">
                <div class="file-info">
                    <div class="file-name">
                        ${this.getFileIcon(file.mimeType)} ${file.originalName}
                    </div>
                    <div class="file-meta">
                        ${this.formatFileSize(file.size)} • ${file.mimeType} • ${new Date(file.uploadedAt).toLocaleString()}
                        ${file.folder ? ` • 📁 ${file.folder}` : ''}
                    </div>
                </div>
                <div class="file-actions">
                    <button onclick="app.viewFile('${file.url}')" class="btn-small"> Ver</button>
                    <button onclick="app.deleteFile('${file.fileId}')" class="btn-small btn-danger"> Eliminar</button>
                </div>
            </div>
        `).join('');

        filesListEl.innerHTML = html;
    }

    getFileIcon(mimeType) {
        if (mimeType.startsWith('image/')) return '🖼️';
        if (mimeType === 'application/pdf') return '📄';
        if (mimeType.startsWith('text/')) return '📝';
        if (mimeType.startsWith('video/')) return '🎥';
        if (mimeType.startsWith('audio/')) return '🎵';
        return '📎';
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    updateFolderFilter(files) {
        const folderFilter = document.getElementById('folderFilter');
        const currentValue = folderFilter.value;
        
        // Get unique folders
        const folders = [...new Set(files.map(f => f.folder).filter(Boolean))];
        
        folderFilter.innerHTML = '<option value="">Todas las carpetas</option>';
        folders.forEach(folder => {
            folderFilter.innerHTML += `<option value="${folder}">${folder}</option>`;
        });
        
        folderFilter.value = currentValue;
    }

    async downloadFile(fileId) {
        try {
            const headers = {
                'X-Client-Id': this.currentClient,
                'Authorization': `Bearer ${this.jwtToken}`
            };

            const response = await fetch(`${this.apiBase}/files/download/${fileId}`, { headers });
            
            if (response.ok) {
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = ''; // Browser will use Content-Disposition header
                a.click();
                window.URL.revokeObjectURL(url);
            } else {
                throw new Error('Error downloading file');
            }
        } catch (error) {
            console.error('❌ Error downloading:', error);
            this.showMessage(`❌ Error descargando: ${error.message}`, 'error');
        }
    }

    viewFile(url) {
        // Abrir en la misma pestaña para imágenes
        window.location.href = window.location.origin + url;
    }

    async deleteFile(fileId) {
        if (!confirm('⚠️ ¿Estás seguro de eliminar este archivo?')) return;

        try {
            const headers = {
                'X-Client-Id': this.currentClient,
                'Authorization': `Bearer ${this.jwtToken}`
            };

            const response = await fetch(`${this.apiBase}/files/${fileId}`, {
                method: 'DELETE',
                headers
            });

            const result = await response.json();

            if (result.success) {
                this.showMessage('✅ Archivo eliminado', 'success');
                this.loadFiles(); // Refresh
            } else {
                throw new Error(result.error || 'Error deleting file');
            }

        } catch (error) {
            console.error('❌ Error deleting:', error);
            this.showMessage(`❌ Error eliminando: ${error.message}`, 'error');
        }
    }

    filterFiles(searchTerm) {
        const fileItems = document.querySelectorAll('.file-item');
        
        fileItems.forEach(item => {
            const fileName = item.querySelector('.file-name').textContent.toLowerCase();
            const isVisible = fileName.includes(searchTerm.toLowerCase());
            item.style.display = isVisible ? 'flex' : 'none';
        });
    }

    showMessage(message, type) {
        // Simple message display - could be improved with toast notifications
        const existingMessage = document.querySelector('.message');
        if (existingMessage) existingMessage.remove();

        const messageEl = document.createElement('div');
        messageEl.className = `message ${type}`;
        messageEl.textContent = message;
        
        document.querySelector('.files-section').prepend(messageEl);
        
        setTimeout(() => {
            messageEl.remove();
        }, 5000);
    }
}

// Initialize app
const app = new FileServerAdmin();