<!--
  ~ Copyright 2025 Simon Emms <simon@simonemms.com>
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

<script lang="ts">
  let { data } = $props();

  let userAccounts = $state(Object.entries(data.user.accounts));

  async function unlinkAccount(providerId: string) {
    await fetch(`/auth/user/provider/${providerId}`, {
      headers: {
        authorization: `Bearer ${data.token}`,
      },
      method: 'delete',
    });

    userAccounts = userAccounts.filter(([provider]) => provider !== providerId);
  }
</script>

This is a logged in page. The user that is logged in is {data.user.name}

<h2>Providers</h2>

<dl>
  {#each userAccounts as [providerId, account]}
    <dt>{providerId}</dt>
    <dd><strong>Provider user ID</strong>: {account.providerUserId}</dd>
    <dd><strong>Name</strong>: {account.name}</dd>
    <dd><strong>Username</strong>: {account.username}</dd>
    <dd><strong>Email address</strong>: {account.emailAddress}</dd>

    {#if userAccounts.length > 1}
      <button onclick={() => unlinkAccount(providerId)}>
        Unlink account
      </button>
    {/if}
  {/each}
</dl>

<ul>
  <li data-sveltekit-preload-data="false">
    <a href="/login">Manage providers</a>
  </li>
  <li data-sveltekit-preload-data="false">
    <a href="/logout">Logout</a>
  </li>
</ul>
