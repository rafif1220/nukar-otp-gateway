<script>
    import { onMount } from 'svelte';
    let logs = $state([]); // <-- Cukup tambahkan $state di sini

    onMount(async () => {
        const res = await fetch('http://localhost:4000/api/history');
        const result = await res.json();
        logs = result.data || [];
    });
</script>

<h1 class="text-3xl font-bold mb-6">History Log</h1>
<div class="overflow-x-auto bg-base-100 rounded-box shadow">
  <table class="table table-zebra w-full">
    <thead>
      <tr>
        <th>Waktu</th>
        <th>Tujuan</th>
        <th>OTP Code</th>
        <th>Status</th>
        <th>Response Barantum</th>
      </tr>
    </thead>
    <tbody>
      {#each logs as log}
        <tr>
          <td>{new Date(log.created_at).toLocaleString('id-ID')}</td>
          <td class="font-mono">{log.to}</td>
          <td class="font-bold text-primary">{log.otp_code}</td>
          <td>
            <span class="badge {log.status === 'SENT' ? 'badge-success text-white' : (log.status === 'FAILED' ? 'badge-error text-white' : 'badge-warning')}">
                {log.status}
            </span>
          </td>
          <td class="max-w-xs truncate text-xs font-mono" title={log.response}>{log.response || '-'}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
