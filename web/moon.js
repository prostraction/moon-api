import CONFIG from './CONFIG.js';
import { createMoonShape, updateMoonShape } from './moonShape.js';

const resultDiv = document.getElementById('result');
const moonHost = document.getElementById('moon');
if (moonHost) createMoonShape(moonHost);

// --- API ---
async function getMoonData(date = new Date()) {
    const params = {
        utc: -Math.round(date.getTimezoneOffset() / 60),
        day: date.getDate(),
        month: date.getMonth() + 1,
        year: date.getFullYear(),
        hour: date.getHours(),
        minute: date.getMinutes(),
        second: date.getSeconds(),
        lang: 'ru'
    };

    const url = new URL(CONFIG.API_URL);
    Object.entries(params).forEach(([key, value]) => url.searchParams.append(key, value.toString()));

    const response = await fetch(url.toString(), {
        method: 'GET',
        headers: { 'Accept': 'application/json' }
    });

    if (!response.ok) throw new Error(`Ошибка API: ${response.status} ${response.statusText}`);
    return await response.json();
}

// --- Отображение ---
export async function showMoonDay(date, isCurrent) {
    const currentHeight = resultDiv.offsetHeight;
    resultDiv.style.minHeight = `${currentHeight}px`;

    try {
        const data = await getMoonData(date);

        let moonDay = Math.floor(data.EndDay.MoonDays);
        let illumination = data.EndDay.Illumination;
        let phase = data.EndDay.Phase;
        let zodiac = data.EndDay.Zodiac;

        if (isCurrent) {
            moonDay = Math.floor(data.CurrentState.MoonDays);
            illumination = data.CurrentState.Illumination;
            phase = data.CurrentState.Phase;
            zodiac = data.CurrentState.Zodiac;
        }

        const formattedDate = date.toLocaleDateString("ru-RU");

        resultDiv.innerHTML = `
            <div class="moon-day">
                <span class="detail-value">${phase.Emoji} Moon day: ${moonDay}</span>
            </div>

            <div class="moon-details">
                <div class="detail-item"><span class="detail-label">Date:</span><span class="detail-value">${formattedDate}</span></div>
                <div class="detail-item"><span class="detail-label">Moon Phase:</span><span class="detail-value">${phase.Name}</span></div>
                <div class="detail-item"><span class="detail-label">Illumination:</span><span class="detail-value">${illumination}%</span></div>
                <div class="detail-item"><span class="detail-label">Zodiac sign:</span><span class="detail-value">${zodiac.Name}</span></div>
            </div>
        `;

        if (moonHost) {
            updateMoonShape(moonHost, illumination / 100, !!phase.IsWaxing);
        }
    } catch (err) {
        resultDiv.innerHTML = `
            <div class="error-title">Ошибка получения данных</div>
            <div class="error-detail">${err.message || 'Неизвестная ошибка'}</div>
        `;
    }
}

export { getMoonData };