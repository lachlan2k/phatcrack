<template>
  <table class="table w-full">
    <!-- head -->
    <thead>
      <tr>
        <th>Project</th>
        <th>Hashlist</th>
        <th>Cracked</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody class="first-col-bold">
      <!-- row 1 -->
      <tr class="hover">
        <td>41234 - Customer A</td>
        <td>System (NTLM)</td>
        <td>25/100 (25%)</td>
        <td><div class="badge-accent badge">Complete</div></td>
      </tr>
      <!-- row 2 -->
      <tr class="hover">
        <td>41235 - Customer B</td>
        <td>Web server (MD5)</td>
        <td>0/300 (0%)</td>
        <td><div class="badge-info badge">Running</div></td>
      </tr>
      <!-- row 3 -->
      <tr class="hover">
        <td>41236 - Customer C</td>
        <td>MSSQL Dump (SHA1)</td>
        <td>10/200 (5%)</td>
        <td><div class="badge-error badge">Failed</div></td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
thead > tr > th {
  background: none !important;
  border-bottom-width: 1px;
  /* border-bottom: 1px solid black; */
}
</style>