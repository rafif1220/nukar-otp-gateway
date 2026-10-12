<script>
    import { onMount } from 'svelte';
    import Chart from 'chart.js/auto';

    let metrics = $state({ sent: 0, failed: 0 }); 
    
    // Variabel untuk Chart.js
    let chartCanvas;
    let chartInstance;
    let selectedRange = $state('7d'); // Default ke 7 Hari Terakhir

    // Fungsi fetch metrics (angka total) dari Golang
    onMount(async () => {
        try {
            const res = await fetch('http://localhost:4000/api/metrics');
            metrics = await res.json();
        } catch (error) {
            console.error("Gagal load metrics:", error);
        }
    });

    // Svelte 5 Rune: Otomatis render ulang grafik tiap kali 'selectedRange' berubah
    $effect(() => {
        if (chartCanvas) {
            renderChart(selectedRange);
        }
    });

    async function renderChart(range) {
        if (chartInstance) {
            chartInstance.destroy();
        }

        try {
            // Tembak API Golang dengan parameter range
            const res = await fetch(`http://localhost:4000/api/chart-data?range=${range}`);
            const realData = await res.json();

            chartInstance = new Chart(chartCanvas, {
                type: 'line',
                data: {
                    labels: realData.labels, // Data label dari Golang
                    datasets: [
                        {
                            label: 'Sukses Terkirim (SENT)',
                            data: realData.sent, // Data sent dari Golang
                            borderColor: '#36D399',
                            backgroundColor: 'rgba(54, 211, 153, 0.1)',
                            borderWidth: 3,
                            tension: 0.4,
                            fill: true,
                            pointRadius: 4
                        },
                        {
                            label: 'Gagal (FAILED)',
                            data: realData.failed, // Data failed dari Golang
                            borderColor: '#F87272',
                            backgroundColor: 'rgba(248, 114, 114, 0.1)',
                            borderWidth: 2,
                            tension: 0.4,
                            fill: true,
                            pointRadius: 3
                        }
                    ]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: { mode: 'index', intersect: false },
                    plugins: { legend: { position: 'top' } },
                    scales: {
                        y: { beginAtZero: true, grid: { borderDash: [5, 5] } },
                        x: { grid: { display: false } }
                    }
                }
            });
        } catch (error) {
            console.error("Gagal render chart dengan data asli:", error);
        }
    }
</script>

<!-- Header & Tombol Refresh -->
<div class="flex justify-between items-center mb-6">
    <h1 class="text-3xl font-bold text-base-content">Dashboard Overview</h1>
    <button class="btn btn-primary btn-sm" onclick={() => location.reload()}>Refresh Data</button>
</div>

<!-- Card Angka Total (Metrics) -->
<div class="card bg-base-100 shadow-sm border-2 border-base-300 mb-6">
    <div class="card-body p-0">
        <div class="stats stats-vertical lg:stats-horizontal w-full bg-transparent">
            <div class="stat border-b lg:border-b-0 lg:border-r border-base-300">
                <div class="stat-figure text-success">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                </div>
                <div class="stat-title font-semibold text-base-content/70">Total OTP Terkirim</div>
                <div class="stat-value text-success">{metrics.sent}</div>
                <div class="stat-desc font-medium">WhatsApp & Email</div>
            </div>
            
            <div class="stat">
                <div class="stat-figure text-error">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                </div>
                <div class="stat-title font-semibold text-base-content/70">Total Gagal</div>
                <div class="stat-value text-error">{metrics.failed}</div>
                <div class="stat-desc text-error font-bold">Butuh eskalasi jika tinggi</div>
            </div>
        </div>
    </div>
</div>

<!-- Card Grafik Chart.js -->
<div class="card bg-base-100 shadow-sm border-2 border-base-300">
    <div class="card-body p-4 sm:p-6">
        <!-- Header Grafik & Dropdown Filter -->
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4 gap-4">
            <h2 class="text-xl font-bold text-base-content">Tren Pengiriman OTP</h2>
            
            <!-- Dropdown Filter Waktu -->
            <select bind:value={selectedRange} class="select select-bordered select-sm w-full sm:w-[120px] font-medium">
                <option value="today">Hari Ini</option>
                <option value="7d">7 Hari Terakhir</option>
                <option value="1m">1 Bulan Terakhir</option>
                <option value="6m">6 Bulan Terakhir</option>
                <option value="1y">1 Tahun Terakhir</option>
            </select>
        </div>

        <!-- Wadah Canvas (tinggi di-set agar responsif) -->
        <div class="relative w-full h-[300px] sm:h-[400px]">
            <canvas bind:this={chartCanvas}></canvas>
        </div>
    </div>
</div>