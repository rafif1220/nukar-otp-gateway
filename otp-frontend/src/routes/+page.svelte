<script>
    import { onMount } from 'svelte';
    let metrics = { sent: 0, failed: 0 };

    onMount(async () => {
        // Fetch ke endpoint Golang lu
        const res = await fetch('http://localhost:4000/api/metrics');
        metrics = await res.json();
    });
</script>

<h1 class="text-3xl font-bold mb-6">Overview</h1>
<div class="stats shadow w-full">
  <div class="stat">
    <div class="stat-title">Total OTP Terkirim</div>
    <div class="stat-value text-success">{metrics.sent}</div>
  </div>
  <div class="stat">
    <div class="stat-title">Total Gagal</div>
    <div class="stat-value text-error">{metrics.failed}</div>
    <div class="stat-desc">Butuh eskalasi jika tinggi</div>
  </div>
</div>
