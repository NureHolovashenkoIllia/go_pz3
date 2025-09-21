// --- CONFIGURATION ---
const API_BASE_URL = '/api/v1';
let currentSection = 'actors';

// --- DOM ELEMENTS ---
// Отримуємо модальне вікно один раз, щоб не шукати його постійно
const crudModalElement = document.getElementById('crudModal');
const crudModal = new bootstrap.Modal(crudModalElement);
const modalTitle = document.getElementById('modalTitle');
const modalForm = document.getElementById('modalForm');
const formFields = document.getElementById('form-fields');
const editIdInput = document.getElementById('editId');

/**
 * Конфігурація для кожної секції UI.
 * Зберігає ендпоїнт, функцію для завантаження даних та поля для модального вікна.
 */
const sectionConfig = {
    'actors': {
        fetchFunc: fetchActors,
        endpoint: 'actors',
        title: 'Актори',
        modalFields: `
            <div class="mb-3"><label class="form-label">ПІБ</label><input type="text" class="form-control" name="FullName" required></div>
            <div class="mb-3"><label class="form-label">Дата народження</label><input type="date" class="form-control" name="BirthDate" required></div>
            <div class="mb-3"><label class="form-label">Стаж (років)</label><input type="number" class="form-control" name="ExperienceYears"></div>
            <div class="mb-3"><label class="form-label">Контактна інформація</label><input type="text" class="form-control" name="ContactInfo"></div>
        `
    },
    'performances': {
        fetchFunc: fetchPerformances,
        endpoint: 'performances',
        title: 'Вистави',
        modalFields: `
            <div class="mb-3"><label class="form-label">Назва</label><input type="text" class="form-control" name="Title" required></div>
            <div class="mb-3"><label class="form-label">Дата прем'єри</label><input type="date" class="form-control" name="PremiereDate"></div>
            <div class="mb-3"><label class="form-label">Жанр</label><input type="text" class="form-control" name="Genre"></div>
            <div class="mb-3"><label class="form-label">Тривалість (хвилин)</label><input type="number" class="form-control" name="DurationMinutes"></div>
        `
    },
    'roles': {
        fetchFunc: fetchRoles,
        endpoint: 'roles',
        title: 'Ролі',
        modalFields: `
            <div class="mb-3"><label class="form-label">Назва ролі</label><input type="text" class="form-control" name="RoleName" required></div>
            <div class="mb-3"><label class="form-label">ID Актора</label><input type="number" class="form-control" name="ActorID" required></div>
            <div class="mb-3"><label class="form-label">ID Вистави</label><input type="number" class="form-control" name="PerformanceID" required></div>
        `
    },
    'rehearsals': {
        fetchFunc: fetchRehearsals,
        endpoint: 'rehearsals',
        title: 'Репетиції',
        modalFields: `
            <div class="mb-3"><label class="form-label">ID Вистави</label><input type="number" class="form-control" name="performance_id" required></div>
            <div class="mb-3"><label class="form-label">Дата і час</label><input type="datetime-local" class="form-control" name="date_time" required></div>
            <div class="mb-3"><label class="form-label">ID Акторів (через кому)</label><input type="text" class="form-control" name="actor_ids" required></div>
        `
    },
    'analytics': {
        fetchFunc: fetchAnalytics,
        title: 'Аналітика'
    }
};


// --- NAVIGATION ---

/**
 * Перемикає видиму секцію на сторінці.
 * @param {string} sectionId - Ідентифікатор секції ('actors', 'performances' тощо).
 */
function switchSection(sectionId) {
    currentSection = sectionId;
    document.querySelectorAll('.content-section').forEach(section => {
        section.classList.toggle('active', section.id === `${sectionId}-section`);
    });

    // Викликаємо функцію завантаження даних для активної секції
    sectionConfig[sectionId].fetchFunc();
}


// --- GENERIC API HELPERS ---

/**
 * Універсальна функція для отримання даних з API.
 * @param {string} endpoint - Ендпоїнт для запиту (напр., 'actors').
 * @returns {Promise<Array|Object>} - Дані з сервера.
 */
async function apiFetch(endpoint) {
    try {
        const response = await fetch(`${API_BASE_URL}/${endpoint}`);
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || `HTTP error! Status: ${response.status}`);
        }
        return response.json();
    } catch (error) {
        console.error(`Error fetching from ${endpoint}:`, error);
        alert(`Помилка завантаження даних: ${error.message}`);
        return null; // Повертаємо null у разі помилки, щоб уникнути падіння
    }
}


// --- MODAL & FORM LOGIC ---

/**
 * Показує модальне вікно для створення нового запису.
 */
function showCreateModal() {
    const config = sectionConfig[currentSection];
    modalTitle.innerText = `Створити: ${config.title}`;
    editIdInput.value = '';
    formFields.innerHTML = config.modalFields;
    modalForm.reset();
    crudModal.show();
}

/**
 * Показує модальне вікно для редагування існуючого запису.
 * @param {number} id - ID запису для редагування.
 */
async function showEditModal(id) {
    const config = sectionConfig[currentSection];
    const data = await apiFetch(`${config.endpoint}/${id}`);
    if (!data) return;

    modalTitle.innerText = `Редагувати: ${config.title} #${id}`;
    editIdInput.value = id;
    formFields.innerHTML = config.modalFields;

    for (const key in data) {
        if (modalForm.elements[key]) {
            const inputElement = modalForm.elements[key];
            let value = data[key];

            if (inputElement.type === 'date') {
                value = new Date(value).toISOString().slice(0, 10);
            } else if (inputElement.type === 'datetime-local') {
                value = formatForDateTimeLocalInput(value);
            }

            inputElement.value = value;
        }
    }

    if (currentSection === 'rehearsals' && data.Actors) {
        modalForm.elements['performance_id'].value = data.PerformanceID;
        modalForm.elements['actor_ids'].value = data.Actors.map(actor => actor.ID).join(', ');

        if (data.DateTime) {
            modalForm.elements['date_time'].value = formatForDateTimeLocalInput(data.DateTime);
        }
    }

    crudModal.show();
}

/**
 * Обробляє відправку форми з модального вікна (створення або оновлення).
 */
async function handleFormSubmit() {
    const config = sectionConfig[currentSection];
    const id = editIdInput.value;
    const formData = new FormData(modalForm);
    let data = Object.fromEntries(formData.entries());

    // Конвертуємо дані у правильні типи перед відправкою
    data = prepareDataForApi(data);

    const method = id ? 'PUT' : 'POST';
    const url = id ? `${API_BASE_URL}/${config.endpoint}/${id}` : `${API_BASE_URL}/${config.endpoint}`;

    try {
        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || JSON.stringify(error));
        }

        crudModal.hide();
        config.fetchFunc(); // Оновлюємо таблицю
    } catch (error) {
        alert(`Помилка збереження: ${error.message}`);
    }
}

/**
 * Видаляє запис за його ID.
 * @param {number} id - ID запису для видалення.
 */
async function deleteItem(id) {
    const config = sectionConfig[currentSection];
    if (confirm(`Ви впевнені, що хочете видалити запис #${id} з "${config.title}"?`)) {
        try {
            const response = await fetch(`${API_BASE_URL}/${config.endpoint}/${id}`, { method: 'DELETE' });
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Помилка на сервері');
            }
            await config.fetchFunc(); // Оновлюємо таблицю
        } catch(error) {
            alert(`Помилка видалення: ${error.message}`);
        }
    }
}

/**
 * Допоміжна функція для конвертації строкових даних з форми у потрібні типи.
 * @param {object} data - Об'єкт з даними форми.
 * @returns {object} - Об'єкт з конвертованими даними.
 */
function prepareDataForApi(data) {
    if (currentSection === 'rehearsals' && data.actor_ids) {
        data.actor_ids = data.actor_ids.split(',').map(id => parseInt(id.trim()));
        data.performance_id = parseInt(data.performance_id);
        data.date_time = new Date(data.date_time).toISOString();
    } else {
        for(let key in data) {
            if (key.endsWith('ID') || key.includes('Years') || key.includes('Minutes')) {
                data[key] = parseInt(data[key]) || 0;
            }
            if ((key === 'BirthDate' || key === 'PremiereDate') && data[key]) {
                data[key] = new Date(data[key]).toISOString();
            }
        }
    }
    return data;
}

/**
 * Правильно форматує дату для поля <input type="datetime-local">.
 * Враховує локальний часовий пояс користувача, а не UTC.
 * @param {Date|string} dateInput - Дата для форматування.
 * @returns {string} - Рядок у форматі YYYY-MM-DDTHH:mm.
 */
function formatForDateTimeLocalInput(dateInput) {
    if (!dateInput) return '';

    const date = new Date(dateInput);

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    return `${year}-${month}-${day}T${hours}:${minutes}`;
}

// --- DATA FETCH & RENDER FUNCTIONS ---

// Актори
async function fetchActors() {
    const data = await apiFetch('actors');
    const tableBody = document.getElementById('actors-table');
    tableBody.innerHTML = data?.map(item => `
        <tr>
            <td>${item.ID}</td>
            <td>${item.FullName}</td>
            <td>${new Date(item.BirthDate).toLocaleDateString()}</td>
            <td>${item.ExperienceYears}</td>
            <td>${item.ContactInfo}</td>
            <td>
                <button class="btn btn-warning btn-sm" onclick="showEditModal(${item.ID})">✏️</button>
                <button class="btn btn-danger btn-sm" onclick="deleteItem(${item.ID})">🗑️</button>
            </td>
        </tr>
    `).join('') || '<tr><td colspan="6">Дані відсутні</td></tr>';
}

// Вистави
async function fetchPerformances() {
    const data = await apiFetch('performances');
    const tableBody = document.getElementById('performances-table');
    tableBody.innerHTML = data?.map(item => `
        <tr>
            <td>${item.ID}</td>
            <td>${item.Title}</td>
            <td>${new Date(item.PremiereDate).toLocaleDateString()}</td>
            <td>${item.Genre}</td>
            <td>${item.DurationMinutes}</td>
            <td>
                <button class="btn btn-warning btn-sm" onclick="showEditModal(${item.ID})">✏️</button>
                <button class="btn btn-danger btn-sm" onclick="deleteItem(${item.ID})">🗑️</button>
            </td>
        </tr>
    `).join('') || '<tr><td colspan="6">Дані відсутні</td></tr>';
}

// Ролі
async function fetchRoles() {
    const data = await apiFetch('roles');
    const tableBody = document.getElementById('roles-table');
    tableBody.innerHTML = data?.map(item => `
        <tr>
            <td>${item.ID}</td>
            <td>${item.RoleName}</td>
            <td>${item.ActorID}</td>
            <td>${item.PerformanceID}</td>
            <td>
                <button class="btn btn-warning btn-sm" onclick="showEditModal(${item.ID})">✏️</button>
                <button class="btn btn-danger btn-sm" onclick="deleteItem(${item.ID})">🗑️</button>
            </td>
        </tr>
    `).join('') || '<tr><td colspan="5">Дані відсутні</td></tr>';
}

// Репетиції
async function fetchRehearsals() {
    const data = await apiFetch('rehearsals');
    const tableBody = document.getElementById('rehearsals-table');
    tableBody.innerHTML = data?.map(item => `
        <tr>
            <td>${item.ID}</td>
            <td>${item.PerformanceID}</td>
            <td>${new Date(item.DateTime).toLocaleString()}</td>
            <td>${item.Actors?.map(a => a.FullName).join(', ') || 'Немає даних'}</td>
            <td>
                <button class="btn btn-warning btn-sm" onclick="showEditModal(${item.ID})">✏️</button>
                <button class="btn btn-danger btn-sm" onclick="deleteItem(${item.ID})">🗑️</button>
            </td>
        </tr>
    `).join('') || '<tr><td colspan="5">Дані відсутні</td></tr>';
}

// Аналітика
async function fetchAnalytics() {
    // Helper function to update analytics card
    const updateCard = (elementId, value, subtext = '') => {
        document.getElementById(elementId).innerText = value || 'Немає даних';
        if (subtext) {
             document.getElementById(elementId.replace('-name', '-count').replace('-title', '-count')).innerText = subtext;
        }
    };

    const mostActive = await apiFetch('analytics/most-active-actor');
    updateCard('most-active-actor-name', mostActive?.full_name, `${mostActive?.role_count || 0} ролей`);

    const leastActive = await apiFetch('analytics/least-active-actor');
    updateCard('least-active-actor-name', leastActive?.full_name, `${leastActive?.role_count || 0} ролей`);

    const perfMostActors = await apiFetch('analytics/performance-most-actors');
    updateCard('performance-most-actors-title', perfMostActors?.title, `${perfMostActors?.actor_count || 0} акторів`);

    const mostFreqPerf = await apiFetch('analytics/most-frequent-performance');
    updateCard('most-frequent-performance-title', mostFreqPerf?.title, `${mostFreqPerf?.rehearsal_count || 0} репетицій`);
}


// --- INITIALIZATION ---

/**
 * Ініціалізує додаток після завантаження DOM.
 */
document.addEventListener('DOMContentLoaded', () => {
    // Налаштовуємо слухачів для навігації
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            switchSection(e.target.dataset.section);
        });
    });

    // Налаштовуємо слухача для кнопки "Зберегти" у модальному вікні
    document.getElementById('saveButton').addEventListener('click', handleFormSubmit);

    // Показуємо першу секцію при завантаженні
    switchSection('actors');
});