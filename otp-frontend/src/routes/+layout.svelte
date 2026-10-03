<script>
    import '../app.css';
</script>

<div class="drawer lg:drawer-open">
  <input id="my-drawer-2" type="checkbox" class="drawer-toggle" />
  <div class="drawer-content flex flex-col p-8 bg-base-200 min-h-screen">
    <!-- Konten halaman akan masuk ke tag slot ini -->
    <slot />
  </div> 
  <div class="drawer-side border-r border-base-300 shadow-sm">
    <label for="my-drawer-2" aria-label="close sidebar" class="drawer-overlay"></label> 
    <ul class="menu p-4 w-64 min-h-full bg-base-100 text-base-content font-semibold space-y-2">
      <li class="menu-title text-primary mb-4 text-lg">NUKAR OTP</li>
      <li><a href="/">📊 Overview</a></li>
      <li><a href="/history">📋 History Log</a></li>
      <li><a href="/message">💬 Manual Message</a></li>
    </ul>
  </div>
</div>
