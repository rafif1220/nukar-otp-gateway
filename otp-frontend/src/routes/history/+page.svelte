<script>
    import { onMount } from 'svelte';
    
    let logs = $state([]);
    
    // State untuk Filter
    let searchTarget = $state('');
    let filterStatus = $state('ALL');
    let filterChannel = $state('ALL');
    let filterDate = $state(''); // State baru untuk filter tanggal

    // Derived state: Otomatis mem-filter data setiap kali input diubah
    let filteredLogs = $derived(logs.filter(log => {
        const matchTarget = log.to.toLowerCase().includes(searchTarget.toLowerCase());
        const matchStatus = filterStatus === 'ALL' || log.status === filterStatus;
        
        const channelStr = log.to.includes('@') ? 'EMAIL' : 'WA';
        const matchChannel = filterChannel === 'ALL' || channelStr === filterChannel;
        
        // Filter Tanggal (mencocokkan string YYYY-MM-DD)
        let matchDate = true;
        if (filterDate) {
            // Gunakan created_at (atau updated_at sebagai fallback)
            const logDateStr = log.created_at || log.updated_at || "";
            matchDate = logDateStr.startsWith(filterDate);
        }
        
        return matchTarget && matchStatus && matchChannel && matchDate;
    }));

    function downloadCSV() {
        if (filteredLogs.length === 0) return;
        
        // Tambahkan header Waktu
        const headers = ['Waktu', 'Tujuan', 'Channel', 'Kode OTP', 'Status', 'Response'];
        const rows = filteredLogs.map(log => [
            log.created_at || log.updated_at || "-", // Data waktu
            log.to,
            log.to.includes('@') ? 'Email' : 'WhatsApp',
            log.otp_code,
            log.status,
            `"${(log.response || '').replace(/"/g, '""')}"`
        ]);
        
        const csvContent = [headers.join(','), ...rows.map(e => e.join(','))].join('\n');
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const url = URL.createObjectURL(blob);
        
        const link = document.createElement("a");
        link.href = url;
        link.download = `nukar_otp_history_${new Date().toISOString().slice(0,10)}.csv`;
        link.click();
        URL.revokeObjectURL(url);
    }

    onMount(async () => {
        try {
            const res = await fetch('http://localhost:4000/api/history');
            const result = await res.json();
            if (result.data) logs = result.data;
        } catch (error) {
            console.error("Gagal mengambil data:", error);
        }
    });
</script>

<div class="flex justify-between items-center mb-6">
    <h1 class="text-3xl font-bold text-base-content">History Log</h1>
    <div class="flex gap-2">
        <button class="btn btn-outline btn-sm" onclick={downloadCSV}>
            Download CSV
        </button>
        <button class="btn btn-primary btn-sm" onclick={() => location.reload()}>Refresh</button>
    </div>
</div>

<!-- Area Filter -->
<div class="card bg-base-100 shadow-sm border border-base-200 mb-6">
    <div class="card-body p-4 flex flex-row flex-wrap gap-4 items-end">
        
        <div class="form-control w-full max-w-[250px]">
            <label class="label"><span class="label-text">Tanggal</span></label>
            <!-- Input date otomatis menghasilkan format YYYY-MM-DD -->
            <input type="date" bind:value={filterDate} class="input input-bordered input-sm" />
        </div>

        <div class="form-control w-full max-w-[250px]">
            <label class="label"><span class="label-text">Cari Nomor / Email</span></label>
            <input type="text" bind:value={searchTarget} placeholder="Ketik nomor/email..." class="input input-bordered input-sm" />
        </div>
        
        <div class="form-control w-full max-w-[150px]">
            <label class="label"><span class="label-text">Status</span></label>
            <select bind:value={filterStatus} class="select select-bordered select-sm w-full">
                <option value="ALL">Semua</option>
                <option value="SENT">SENT</option>
                <option value="FAILED">FAILED</option>
                <option value="PENDING">PENDING</option>
            </select>
        </div>

        <div class="form-control w-full max-w-[150px]">
            <label class="label"><span class="label-text">Channel</span></label>
            <select bind:value={filterChannel} class="select select-bordered select-sm w-full">
                <option value="ALL">Semua</option>
                <option value="WA">WhatsApp</option>
                <option value="EMAIL">Email</option>
            </select>
        </div>
        
        <!-- Tombol Reset Filter -->
        <button class="btn btn-ghost btn-sm" onclick={() => { filterDate = ''; searchTarget = ''; filterStatus = 'ALL'; filterChannel = 'ALL'; }}>
            Reset
        </button>
    </div>
</div>

<!-- Area Tabel -->
<div class="card bg-base-100 shadow-sm border-2 border-base-200">
    <div class="card-body p-0 overflow-x-auto">
        <!-- Tambahkan border-2 pada tabel -->
        <table class="table table-zebra w-full border-2 border-base-200">
            
            <!-- Header dengan garis bawah yang tebal -->
            <thead class="bg-base-100 border-b-2 border-base-200 text-base-content">
                <tr>
                    <th class="border-r border-base-200">Waktu</th>
                    <th class="border-r border-base-200">Tujuan (To)</th>
                    <th class="border-r border-base-200">Channel</th>
                    <th class="border-r border-base-200">Kode OTP</th>
                    <th class="border-r border-base-200">Status</th>
                    <th>Response Barantum/SMTP</th>
                </tr>
            </thead>
            
            <tbody>
                {#if filteredLogs.length === 0}
                    <tr>
                        <td colspan="7" class="text-center py-8 text-gray-500">Data tidak ditemukan.</td>
                    </tr>
                {:else}
                    {#each filteredLogs as log}
                        <!-- Setiap baris diberi garis bawah tebal -->
                        <tr class="border-b-2 border-base-300 hover">
                            <td class="text-xs whitespace-nowrap border-r border-base-200">{(log.created_at || log.updated_at || "-").substring(0, 19).replace('T', ' ')}</td>
                            <td class="font-medium border-r border-base-200">{log.to}</td>
                            <td class="border-r border-base-200">
                                {#if log.to.includes('@')}
                                    <span class="badge badge-info badge-sm text-white border-none">Email</span>
                                {:else}
                                    <span class="badge badge-success badge-sm text-white border-none">WhatsApp</span>
                                {/if}
                            </td>
                            <td class="font-mono bg-base-200 text-center px-2 border-r border-base-300 font-bold">{log.otp_code}</td>
                            <td class="border-r border-base-200">
                                {#if log.status === 'SENT'}
                                    <div class="badge badge-success text-white border-none shadow-sm">SENT</div>
                                {:else if log.status === 'FAILED'}
                                    <div class="badge badge-error text-white border-none shadow-sm">FAILED</div>
                                {:else}
                                    <div class="badge badge-warning text-white border-none shadow-sm">PENDING</div>
                                {/if}
                            </td>
                            <td>
                                <div class="max-w-md max-h-24 overflow-y-auto bg-base-200/50 p-2 rounded text-xs font-mono whitespace-pre-wrap select-all border border-base-300 shadow-inner">
                                    {log.response}
                                </div>
                            </td>
                        </tr>
                    {/each}
                {/if}
            </tbody>
        </table>
    </div>
</div>