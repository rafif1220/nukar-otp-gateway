<script>
    let phone = '';
    let message = '';
    let notif = '';

    async function sendMessage() {
        if (!phone || !message) return;
        
        const res = await fetch('http://localhost:4000/api/send-message', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ phone, message })
        });
        const data = await res.json();
        
        notif = data.message;
        phone = ''; 
        message = '';
        
        // Hilangkan notif setelah 3 detik
        setTimeout(() => notif = '', 3000);
    }
</script>

<h1 class="text-3xl font-bold mb-6">Kirim Pesan Manual</h1>
<div class="card bg-base-100 shadow-xl max-w-md">
    <div class="card-body">
        {#if notif} 
            <div class="alert alert-success p-2 mb-4 text-white text-sm">
                <span>{notif}</span>
            </div> 
        {/if}
        
        <div class="form-control">
            <label class="label"><span class="label-text font-semibold">Nomor Tujuan</span></label>
            <input type="text" bind:value={phone} placeholder="Contoh: 08123456789" class="input input-bordered" />
        </div>
        
        <div class="form-control mt-4">
            <label class="label"><span class="label-text font-semibold">Isi Pesan</span></label>
            <textarea bind:value={message} class="textarea textarea-bordered h-24" placeholder="Ketik pesan di sini..."></textarea>
        </div>
        
        <div class="form-control mt-6">
            <button class="btn btn-primary" on:click={sendMessage}>Kirim Pesan</button>
        </div>
    </div>
</div>
