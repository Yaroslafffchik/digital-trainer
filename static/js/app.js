// Подключение к WebSocket
const ws = new WebSocket('ws://localhost:1236/ws');

ws.onopen = () => {
    console.log('Подключено к WebSocket серверу');
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.status) {
        alert(data.status);
    }
};

ws.onclose = () => {
    console.log('WebSocket соединение закрыто');
};

// Загрузка тренировок
function loadTrainings() {
    fetch('/api/trainings')
        .then(response => response.json())
        .then(data => {
            const trainingList = document.getElementById('trainingList');
            if (trainingList) {
                trainingList.innerHTML = '';
                data.forEach(training => {
                    const li = document.createElement('li');
                    li.className = 'list-group-item d-flex justify-content-between align-items-center';
                    li.innerHTML = `
                        <span>${training.name} (Сессия: ${training.session_id})</span>
                        <button class="btn btn-info btn-sm open-description" data-description="${training.description}" data-bs-toggle="modal" data-bs-target="#descriptionModal">Открыть описание</button>
                    `;
                    trainingList.appendChild(li);
                });

                // Обработчик для открытия описания
                document.querySelectorAll('.open-description').forEach(button => {
                    button.addEventListener('click', () => {
                        const description = button.getAttribute('data-description');
                        document.getElementById('modalDescription').textContent = description;
                    });
                });
            }
        });
}

// Загрузка челленджей
function loadChallenges() {
    fetch('/api/challenges')
        .then(response => response.json())
        .then(data => {
            const challengeList = document.getElementById('challengeList');
            if (challengeList) {
                challengeList.innerHTML = '';
                data.forEach(challenge => {
                    const li = document.createElement('li');
                    li.className = 'list-group-item';
                    li.textContent = `${challenge.name}: ${challenge.description} (Создано: ${challenge.created_at})`;
                    challengeList.appendChild(li);
                });
            }
        });
}

// Загрузка метрик и обновление графика
function loadMetrics() {
    fetch('/api/metrics')
        .then(response => response.json())
        .then(data => {
            const metricsList = document.getElementById('metricsList');
            if (metricsList) {
                metricsList.innerHTML = '';
                data.forEach(metric => {
                    const li = document.createElement('li');
                    li.className = 'list-group-item';
                    li.textContent = `Сессия: ${metric.session_id}, Пульс: ${metric.pulse}, Скорость: ${metric.speed}, Время: ${metric.timestamp}`;
                    metricsList.appendChild(li);
                });
            }

            const metricsChart = document.getElementById('metricsChart');
            if (metricsChart) {
                new Chart(metricsChart, {
                    type: 'line',
                    data: {
                        labels: data.map(m => m.timestamp),
                        datasets: [
                            {
                                label: 'Пульс (уд/мин)',
                                data: data.map(m => m.pulse),
                                borderColor: '#007bff',
                                backgroundColor: 'rgba(0, 123, 255, 0.1)',
                                fill: true,
                            },
                            {
                                label: 'Скорость (км/ч)',
                                data: data.map(m => m.speed),
                                borderColor: '#dc3545',
                                backgroundColor: 'rgba(220, 53, 69, 0.1)',
                                fill: true,
                            },
                        ],
                    },
                    options: {
                        responsive: true,
                        scales: {
                            y: {
                                beginAtZero: true,
                            },
                        },
                    },
                });
            }
        });
}

// Отправка тренировки через HTTP
if (document.getElementById('trainingForm')) {
    document.getElementById('trainingForm').addEventListener('submit', (e) => {
        e.preventDefault();
        const sessionId = document.getElementById('sessionId').value;
        const name = document.getElementById('trainingName').value;
        const description = document.getElementById('trainingDescription').value;

        fetch('/api/trainings', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, description, session_id: sessionId }),
        })
            .then(response => response.text())
            .then(data => {
                alert(data);
                loadTrainings();
            })
            .catch(err => console.error('Ошибка:', err));
    });
}


/// Отправка метрик через HTTP POST
function sendMetrics(sessionId, pulse, speed) {
    fetch('/api/metrics', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ session_id: sessionId, pulse: parseInt(pulse), speed: parseFloat(speed) }),
    })
        .then(response => response.text())
        .then(data => {
            alert(data);
            loadMetrics(); // Обновляем метрики только после успешной отправки
        })
        .catch(err => console.error('Ошибка:', err));
}

// Обработка формы метрик
if (document.getElementById('metricsForm')) {
    document.getElementById('metricsForm').addEventListener('submit', (e) => {
        e.preventDefault();
        const sessionId = document.getElementById('sessionIdMetric').value;
        const pulse = document.getElementById('pulse').value;
        const speed = document.getElementById('speed').value;

        sendMetrics(sessionId, pulse, speed);
    });
}

// Периодическое обновление данных
if (document.getElementById('trainingList')) {
    loadTrainings();
    setInterval(loadTrainings, 5000);
}
if (document.getElementById('challengeList')) {
    loadChallenges();
    setInterval(loadChallenges, 5000);
}
if (document.getElementById('metricsList') || document.getElementById('metricsChart')) {
    loadMetrics(); // Начальная загрузка метрик
    setInterval(loadMetrics, 2000); // Обновление каждые 2 секунды
}

