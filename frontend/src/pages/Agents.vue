<template>
  <main class="w-full p-5">
    <div class="prose">
      <h1>Agents</h1>
    </div>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="grow">
        <div class="stats shadow">
          <div class="stat">
            <div class="stat-title">Agents Online</div>
            <div class="stat-value flex justify-between">
              <span>2/3</span>
              <span class="mt-1 text-2xl text-primary">
                <font-awesome-icon icon="fa-solid fa-robot" />
              </span>
            </div>
          </div>

          <div class="stat">
            <div class="stat-title">Total Power Draw</div>
            <div class="stat-value flex justify-between">
              <span>3,801w</span>
              <span class="mt-1 ml-4 text-2xl text-yellow-400">
                <font-awesome-icon icon="fa-solid fa-bolt" />
              </span>
            </div>
          </div>

          <div class="stat">
            <div class="stat-title">Total Attacks Running</div>

            <div class="stat-value flex justify-between">
              <span>8</span>
              <span class="mt-1 text-2xl text-info">
                <font-awesome-icon icon="fa-solid fa-bars-progress" />
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="flex basis-full"></div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Agent List</h2>
          <table class="compact-table table w-full">
            <!-- head -->
            <thead>
              <tr>
                <th>Hostname</th>
                <th>CPU</th>
                <th>GPU</th>
                <th>Power Usage</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody class="first-col-bold">
              <tr class="hover">
                <td>jimbob.lan</td>
                <td>
                  <font-awesome-icon icon="fa-solid fa-microchip" /> Intel Core i7 7700k (93°c)
                </td>
                <td>
                  <font-awesome-icon icon="fa-solid fa-memory" /> RTX 3080 Ti (78°c)
                  <br />
                  <font-awesome-icon icon="fa-solid fa-memory" /> RTX 3080 Ti (81°c)
                  <br />
                  <font-awesome-icon icon="fa-solid fa-memory" /> RTX 3080 Ti (76°c)
                  <br />
                </td>
                <td>1,928w</td>
                <td class="text-center">
                  <div class="badge badge-accent badge-sm"></div>
                </td>
              </tr>

              <tr class="hover">
                <td>workhorse.lan</td>
                <td>
                  <font-awesome-icon icon="fa-solid fa-microchip" /> AMD Radeon Ryzen 7 5800X (83°c)
                </td>
                <td><font-awesome-icon icon="fa-solid fa-memory" /> RTX 3070 Ti (74°c)</td>
                <td>704w</td>
                <td class="text-center">
                  <div class="badge badge-accent badge-sm"></div>
                </td>
              </tr>

              <tr class="hover">
                <td>pizzabox.lan</td>
                <td>
                  <font-awesome-icon icon="fa-solid fa-microchip" /> Intel Xeon 2660 v4
                  <br />
                  <font-awesome-icon icon="fa-solid fa-microchip" /> Intel Xeon 2660 v4
                </td>
                <td><font-awesome-icon icon="fa-solid fa-memory" /> Quadro A5000</td>
                <td>-</td>

                <td class="text-center">
                  <div class="badge badge-ghost badge-sm"></div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
thead > tr > th {
  background: none !important;
  border-bottom-width: 1px;
  /* border-bottom: 1px solid black; */
}

.first-col-bold > tr td:first-of-type {
  font-weight: bold;
}
</style>
