<script>
    import '../app.css';
    
    // State untuk melacak apakah sidebar terbuka penuh atau mode mini
    let isSidebarOpen = $state(true);

    // Fungsi untuk toggle sidebar
    function toggleSidebar() {
        isSidebarOpen = !isSidebarOpen;
    }
</script>

<div class="flex h-screen bg-base-300 overflow-hidden">
    <!-- Sidebar -->
    <!-- Transisi halus untuk perubahan lebar (w-64 ke w-20) -->
    <aside class="bg-base-100 border-r-2 border-base-300 shadow-lg transition-all duration-300 ease-in-out flex flex-col z-20 {isSidebarOpen ? 'w-64' : 'w-20'}">
        
        <!-- Header Sidebar & Tombol Hamburger -->
        <div class="flex items-center justify-between h-16 px-4 border-b border-base-300">
            {#if isSidebarOpen}
                <span class="text-2xl font-black text-primary tracking-wider whitespace-nowrap">OTP-NUKAR</span>
            {/if}
            
            <button onclick={toggleSidebar} class="btn btn-square btn-ghost btn-sm {isSidebarOpen ? '' : 'mx-auto'}" title={isSidebarOpen ? 'Tutup sidebar' : 'Buka sidebar'}>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-6 h-6 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            </button>
        </div>

        <!-- Menu Items -->
        <ul class="menu p-2 flex-1 w-full text-base-content overflow-x-hidden gap-1 mt-2">
            
            <!-- Menu: Overview -->
            <li>
                <!-- Gunakan tooltip DaisyUI agar teks muncul saat mode mini di-hover -->
                <a href="/" class="flex items-center px-4 py-3 {isSidebarOpen ? '' : 'justify-center tooltip tooltip-right'}" data-tip={isSidebarOpen ? '' : 'Overview'}>
                    <!-- Ikon Dashboard/Overview -->
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>
                    
                    {#if isSidebarOpen}
                        <span class="ml-3 font-medium whitespace-nowrap">Overview</span>
                    {/if}
                </a>
            </li>

            <!-- Menu: History Log -->
            <li>
                <a href="/history" class="flex items-center px-4 py-3 {isSidebarOpen ? '' : 'justify-center tooltip tooltip-right'}" data-tip={isSidebarOpen ? '' : 'History Log'}>
                    <!-- Ikon Document/History -->
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                    
                    {#if isSidebarOpen}
                        <span class="ml-3 font-medium whitespace-nowrap">History Log</span>
                    {/if}
                </a>
            </li>

        </ul>
    </aside>

    <!-- Area Konten Utama -->
    <main class="flex-1 overflow-y-auto p-6 transition-all duration-300">
        <slot />
    </main>
</div>