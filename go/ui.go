package main

const gatewayContentSecurityPolicy = "default-src 'none'; script-src 'nonce-gateway-ui'; style-src 'unsafe-inline'; connect-src 'self'; img-src data:; base-uri 'none'; form-action 'none'"

func gatewayUIHTML() string {
	return `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta http-equiv="Content-Security-Policy" content="` + gatewayContentSecurityPolicy + `">
<title>Gateway Manager</title>
<style>
:root{--bg:#f3efe7;--panel:#fffaf2;--ink:#18262a;--muted:#68767a;--line:#d8cfbf;--accent:#0f766e;--warm:#a16207;--danger:#b42318}
*{box-sizing:border-box}body{margin:0;font-family:"Segoe UI",Tahoma,sans-serif;background:radial-gradient(circle at top left,#fff6e3 0,#f3efe7 52%,#eae4d8 100%);color:var(--ink)}
.wrap{max-width:1320px;margin:0 auto;padding:26px 18px 50px}.hero{display:grid;gap:8px;margin-bottom:18px}.hero h1{margin:0;font-size:34px}.hero p{margin:0;color:var(--muted);line-height:1.6;max-width:900px}
.layout{display:grid;grid-template-columns:340px 1fr;gap:18px}.stack{display:grid;gap:16px}.card{background:linear-gradient(180deg,#fffdf9,#fff7ee);border:1px solid var(--line);border-radius:22px;padding:18px;box-shadow:0 12px 36px rgba(72,55,27,.08)}
.card h2{margin:0 0 12px;font-size:18px}.card h3{margin:0 0 10px;font-size:14px;text-transform:uppercase;letter-spacing:.06em;color:var(--muted)}
label{display:grid;gap:6px;font-size:12px;text-transform:uppercase;letter-spacing:.06em;color:var(--muted)} input,textarea,select{width:100%;padding:10px 12px;border:1px solid #cec4b3;border-radius:12px;background:#fff;color:var(--ink);font:inherit}
textarea{min-height:110px;resize:vertical;font-family:ui-monospace,Consolas,monospace}.big{min-height:260px} .grid2{display:grid;grid-template-columns:1fr 1fr;gap:12px}.grid3{display:grid;grid-template-columns:repeat(3,1fr);gap:12px}.grid4{display:grid;grid-template-columns:repeat(4,1fr);gap:12px}
.actions{display:flex;flex-wrap:wrap;gap:10px}button{border:0;border-radius:999px;padding:11px 16px;background:var(--ink);color:#fff;font-weight:600;cursor:pointer}button.alt{background:var(--accent)}button.warn{background:var(--warm)}button.danger{background:var(--danger)}
.list{display:grid;gap:10px}.item{padding:12px;border:1px solid #ddd0bc;border-radius:16px;background:#fffdf9;cursor:pointer}.item.active{border-color:var(--accent);box-shadow:0 0 0 2px rgba(15,118,110,.16)}.pill{display:inline-flex;align-items:center;border-radius:999px;padding:4px 10px;background:#efe7da;color:#5d4d35;font-size:12px}.hint{font-size:12px;color:var(--muted)}
pre{margin:0;padding:14px;border-radius:16px;background:#152126;color:#d9f6f2;overflow:auto;max-height:420px}.stats{display:grid;gap:10px}.stat{padding:12px;border:1px solid #ddd0bc;border-radius:16px;background:#fff}.stat b{display:block;font-size:24px;margin-top:4px}
.section-title{display:flex;align-items:center;justify-content:space-between;gap:12px;margin:10px 0 8px}.section-title h3{margin:0}.rule-builder,.group-board,.route-board{display:grid;gap:12px}.mini-card{padding:12px;border:1px solid #ddd0bc;border-radius:16px;background:#fffdf9}.mini-card h4{margin:0 0 8px;font-size:12px;text-transform:uppercase;letter-spacing:.06em;color:var(--muted)}.chips{display:flex;flex-wrap:wrap;gap:8px}.chip{display:inline-flex;align-items:center;border-radius:999px;padding:6px 10px;background:#efe7da;color:#5d4d35;font-size:12px}.rowline{display:flex;flex-wrap:wrap;gap:8px}.split{display:grid;grid-template-columns:1.2fr .8fr;gap:12px}.mono{font-family:ui-monospace,Consolas,monospace}.soft{background:#f7f1e7}.tablelist{display:grid;gap:8px}.tableitem{display:grid;grid-template-columns:1.2fr .7fr auto;gap:10px;align-items:center;padding:10px 12px;border:1px solid #ddd0bc;border-radius:14px;background:#fff}.tableitem b{font-size:13px}.tableitem span{font-size:12px;color:var(--muted)}.topology-flow{display:grid;gap:10px}.topology-policy{display:grid;gap:10px;padding:12px;border:1px solid #d7ccbb;border-radius:18px;background:linear-gradient(180deg,#fffdf9,#f9f3e7)}.topology-policy.active{box-shadow:0 0 0 3px rgba(15,118,110,.22);border-color:#0f766e}.topology-stage{display:grid;gap:8px;padding-left:18px;border-left:2px solid #d6c8af}.topology-stage.active{border-left-color:#0f766e;background:rgba(15,118,110,.04);border-radius:14px;padding:8px 0 8px 18px}.topology-rule{display:grid;gap:6px;padding-left:20px;position:relative}.topology-rule.active .item{border-color:#0f766e;box-shadow:0 0 0 3px rgba(15,118,110,.18);background:#fffdf8}.topology-rule.preview-focus .item{border-color:#a16207;box-shadow:0 0 0 3px rgba(161,98,7,.18);background:#fff8ef}.topology-rule.preview-focus .node-detail{border-color:#a16207;background:#fff8ef}.topology-rule::before{content:"";position:absolute;left:6px;top:8px;width:8px;height:8px;border-radius:999px;background:#0f766e}.flow-arrow{font-size:11px;color:var(--muted);padding-left:6px}.node-detail{margin-top:8px;padding:10px 12px;border:1px dashed #d7ccbb;border-radius:14px;background:linear-gradient(180deg,#fffdf9,#fff7ee);box-shadow:inset 0 1px 0 rgba(255,255,255,.55)}.node-detail b{display:block;font-size:12px;color:var(--ink);margin-bottom:6px}.node-detail span{display:block;font-size:11px;color:var(--muted);line-height:1.5}.node-detail .actions{margin-top:8px}.subnode-list{display:grid;gap:6px;margin-top:8px}.subnode{display:flex;align-items:center;justify-content:space-between;gap:8px;padding:8px 10px;border:1px solid #ddd0bc;border-radius:12px;background:#fff}.subnode strong{font-size:12px}.subnode .hint{display:block}.subnode.focused{border-color:#a16207;box-shadow:0 0 0 2px rgba(161,98,7,.18);background:#fff7e8}.subnode.route{border-left:4px solid #0f766e}.subnode.failover{border-left:4px solid #a16207}.subnode.mirror{border-left:4px solid #2563eb}
@media (max-width:1050px){.layout,.grid2,.grid3,.grid4,.split{grid-template-columns:1fr}.hero h1{font-size:28px}}
.wrap{max-width:1480px;padding:28px 22px 56px}.hero{position:relative;overflow:hidden;margin-bottom:22px;padding:24px;border:1px solid rgba(24,38,42,.1);border-radius:32px;background:linear-gradient(135deg,rgba(255,250,242,.92),rgba(236,246,242,.72));box-shadow:0 24px 70px rgba(72,55,27,.10)}.hero::after{content:"";position:absolute;right:-80px;top:-90px;width:300px;height:300px;border-radius:999px;background:radial-gradient(circle,rgba(15,118,110,.20),rgba(15,118,110,0) 68%);pointer-events:none}.hero h1{letter-spacing:-.04em}.hero>.card{position:relative;z-index:1;background:rgba(255,253,249,.82);backdrop-filter:blur(12px);box-shadow:none}.layout.tab-shell{display:grid;grid-template-columns:1fr;gap:18px}.tab-nav{position:relative;z-index:1;display:flex;flex-wrap:wrap;gap:8px;padding:8px;border:1px solid rgba(24,38,42,.10);border-radius:24px;background:rgba(255,253,249,.72);box-shadow:0 16px 44px rgba(72,55,27,.08);backdrop-filter:blur(14px)}.tab-button{display:grid;gap:2px;align-content:center;min-height:58px;padding:10px 16px;border:1px solid transparent;border-radius:18px;background:transparent;color:var(--ink);text-align:left}.tab-button b{font-size:14px;letter-spacing:-.01em}.tab-button span{font-size:11px;color:var(--muted);font-weight:500}.tab-button.active{background:#122225;color:#fff;box-shadow:0 12px 28px rgba(18,34,37,.18)}.tab-button.active span{color:#bdd4cf}.tab-content{display:grid;gap:18px}.tab-pane{display:none;grid-template-columns:minmax(0,1fr);gap:18px}.tab-pane.active{display:grid}.tab-pane>.card{position:static!important;top:auto!important;z-index:auto!important}.tab-pane[data-columns="2"].active{grid-template-columns:minmax(280px,.8fr) minmax(0,1.4fr);align-items:start}.tab-pane[data-columns="audit"].active{grid-template-columns:minmax(260px,.75fr) minmax(0,1.25fr);align-items:start}.card{border-radius:24px;background:rgba(255,253,249,.92);box-shadow:0 18px 48px rgba(72,55,27,.08)}.card h2{display:flex;align-items:center;gap:10px;margin-bottom:14px;font-size:18px;letter-spacing:-.02em}.card h2::before{content:"";width:10px;height:10px;border-radius:999px;background:var(--accent);box-shadow:0 0 0 6px rgba(15,118,110,.10)}.actions{align-items:center;gap:8px;margin:12px 0;padding:8px;border:1px solid rgba(24,38,42,.07);border-radius:18px;background:rgba(239,231,218,.38)}button{display:inline-flex;align-items:center;justify-content:center;min-height:38px;padding:9px 14px;border:1px solid rgba(24,38,42,.08);box-shadow:0 7px 18px rgba(24,38,42,.08);white-space:nowrap;transition:transform .14s ease,box-shadow .14s ease,background .14s ease}button:hover{transform:translateY(-1px);box-shadow:0 11px 24px rgba(24,38,42,.12)}button.alt{background:#0f766e}button.warn{background:#b66b0d}button.danger{background:#b42318}.grid2,.grid3,.grid4,.split{align-items:start}label{min-width:0}.item,.mini-card,.stat,.tableitem,.subnode{box-shadow:0 8px 22px rgba(72,55,27,.04)}pre{max-height:360px;border:1px solid rgba(255,255,255,.08);box-shadow:inset 0 1px 0 rgba(255,255,255,.06)}#currentNodeDetail{position:static!important}.management-card{display:grid;gap:10px}.management-row{display:grid;grid-template-columns:minmax(260px,1fr) minmax(180px,.42fr);gap:12px}.context-card{display:flex;flex-wrap:wrap;align-items:center;justify-content:space-between;gap:10px}.context-card strong{font-size:13px;text-transform:uppercase;letter-spacing:.08em;color:var(--muted)}@media (max-width:1050px){.wrap{padding:18px 12px 42px}.hero{padding:18px;border-radius:24px}.tab-nav{overflow-x:auto;flex-wrap:nowrap}.tab-button{min-width:170px}.tab-pane[data-columns="2"].active,.tab-pane[data-columns="audit"].active,.management-row{grid-template-columns:1fr}.actions{overflow-x:auto;flex-wrap:nowrap;justify-content:flex-start}}
</style>
</head>
<body>
<div class="wrap">
  <section class="hero">
    <h1>Gateway Manager</h1>
    <p>Manage per-key gateway routing, model rewrite rules, daily limits, minute rate limits, and dry-run checks for CPA top-level API keys.</p>
    <div class="card management-card" style="padding:14px 16px;border-radius:18px">
      <h2 style="margin:0 0 10px">Management Access</h2>
      <p class="hint" style="margin-top:0">Enter the CPAMC management key before connecting. The key is kept in this browser session only and is sent as management authentication headers.</p>
      <div class="management-row">
        <label>CPAMC Management Key<input id="managementKeyInput" type="password" autocomplete="current-password" placeholder="management secret-key"></label>
        <label>Status<input id="managementAuthStatus" readonly value="Not connected"></label>
      </div>
      <div class="actions"><button class="alt" id="connectManagementBtn">Connect / Refresh</button><button id="saveManagementKeyBtn">Save For Session</button><button class="danger" id="clearManagementKeyBtn">Clear Key</button></div>
    </div>
    <div class="card context-card" style="padding:12px 16px;border-radius:18px"><strong>Current Context</strong><div id="contextBreadcrumb" class="hint">window 1h | no policy selected</div></div>
  </section>
  <section class="layout">
    <aside class="stack">
      <article class="card">
        <h2>Configured Keys</h2>
        <div class="actions"><button class="alt" id="refreshBtn">Refresh</button><button id="newPolicyBtn">Add Key Policy</button><button class="warn" id="clonePolicyBtn">Clone Policy</button><button class="warn" id="exportAuditBtn">Export Audit</button></div>
        <div class="list" id="keysList"></div>
      </article>
      <article class="card">
        <h2>Usage Snapshot</h2>
        <div class="actions"><button class="alt" id="statsWindow5mBtn">5m</button><button class="alt" id="statsWindow1hBtn">1h</button><button class="alt" id="statsWindow24hBtn">24h</button></div>
        <div class="stats" id="usageStats"></div><div class="card" style="margin-top:12px;padding:12px;border-radius:16px"><h3>Pool Operations</h3><div class="hint">Apply fast traffic-management actions to the current route pool.</div><div class="grid2"><label>Canary Secondary<input id="poolCanarySecondaryInput" placeholder="codex/gpt-5.4"></label><label>Canary Split<input id="poolCanaryPercentInput" type="number" min="1" max="99" value="10"></label></div><div class="grid2"><label>Shift Provider<input id="poolShiftProviderInput" placeholder="openai"></label><label>Shift Percent<input id="poolShiftPercentInput" type="number" min="1" max="100" value="50"></label></div><div class="grid2"><label>Safe Apply<select id="poolSafeApplyInput"><option value="false">false</option><option value="true">true</option></select></label><label>Preview Token<input id="poolPreviewTokenInput" readonly placeholder="preview required when safe mode is on"></label></div><div class="actions"><button class="alt" id="poolPreviewCanaryBtn">Preview Canary</button><button class="alt" id="poolPreviewShiftBtn">Preview Shift</button><button class="alt" id="poolPreviewDrainBtn">Preview Drain</button><button class="alt" id="poolPreviewResumeBtn">Preview Resume</button><button class="alt" id="poolPreviewRebalanceBtn">Preview Rebalance</button><button class="alt" id="poolPreviewRestoreBtn">Preview Restore</button></div><div class="actions"><button class="alt" id="poolDrainBtn">Apply Drain</button><button class="alt" id="poolResumeBtn">Apply Resume</button><button class="warn" id="poolCanaryBtn">Apply Canary</button><button class="alt" id="poolRebalanceBtn">Apply Rebalance</button><button class="alt" id="poolRestoreBtn">Apply Restore</button><button class="warn" id="poolShiftProviderBtn">Apply Shift</button></div><pre id="poolPreviewOut">{}</pre></div><div class="card" style="margin-top:12px;padding:12px;border-radius:16px"><h3>Global Topology</h3><div class="hint">Legend: <span class="pill">medium activity</span> <span class="warn">high activity</span> <span class="hint">low activity</span></div><div class="hint" id="legendExplain">High activity means the node is receiving materially more recent hits in the selected window.</div><div id="globalTopologyOut" class="list"></div></div>
      </article>
      <article class="card">
        <h2>Diagnostics</h2>
        <div class="stats" id="healthStats"></div>
        <div class="list" id="healthWarnings"></div>
      </article>
      <article class="card">
        <h2>Activity</h2>
        <pre id="logBox">Ready.</pre>
      </article><article class="card">
        <h2>Guide</h2>
        <div class="hint">warn means high recent activity. pill means medium recent activity. hint means low recent activity.</div><div class="hint">Double-click topology policy/stage nodes to collapse or expand them.</div><div class="hint">URL state remembers current window, key, rule, and collapsed topology.</div>
      </article>
      <article class="card">
        <h2>Audit Summary</h2>
        <div class="stats" id="auditSummaryStats"></div>
      </article><article class="card">
        <h2>Audit Log</h2>
        <div class="grid2"><label>Decision Filter<select id="auditDecisionFilter"><option value="">all</option><option value="rewrite">rewrite</option><option value="reject">reject</option><option value="pass">pass</option></select></label><label>Rule Filter<input id="auditRuleFilter" placeholder="rule id"></label></div><div class="grid2"><label>Reason Filter<input id="auditReasonFilter" placeholder="reason"></label><label>Key Filter<input id="auditKeyFilter" placeholder="masked key"></label></div><div class="grid3"><label>Policy Filter<input id="auditPolicyFilter" placeholder="policy name"></label><label>Model Filter<input id="auditModelFilter" placeholder="gpt-5.4"></label><label>Provider Filter<input id="auditProviderFilter" placeholder="openai"></label></div><div class="grid3"><label>Event Type<input id="auditEventTypeFilter" placeholder="operator"></label><label>Operator<input id="auditOperatorFilter" placeholder="canary-split"></label><label>Member<input id="auditMemberFilter" placeholder="openai/gpt-5.4"></label></div><div class="grid2"><label>From Time<input id="auditFromFilter" placeholder="2026-07-04T00:00:00+08:00"></label><label>To Time<input id="auditToFilter" placeholder="2026-07-04T23:59:59+08:00"></label></div><div class="actions"><button class="alt" id="applyAuditFilterBtn">Apply Audit Filter</button><button class="alt" id="lastHourAuditBtn">Last Hour</button><button class="alt" id="todayAuditBtn">Today</button></div><div class="list" id="auditList"></div>
      </article>
    </aside>
    <main class="stack">
      <article class="card">
        <h2>Selected Policy</h2>
        <div class="grid2">
          <label>Key ID<input id="policyKeyId" readonly></label>
          <label>Display Name<input id="displayName"></label>
        </div>
        <div class="grid2">
          <label>Match API Key<input id="matchApiKey" placeholder="only stored server-side"></label>
          <label>Enabled<select id="policyEnabled"><option value="true">true</option><option value="false">false</option></select></label>
        </div>
        <div class="grid3">
          <label>Requests / Day<input id="requestsPerDay" type="number" min="0"></label>
          <label>Requests / Min<input id="requestsPerMin" type="number" min="0"></label>
          <label>Max Inflight<input id="maxInflight" type="number" min="0"></label>
        </div>
        <div class="grid2">
          <label>Not Before<input id="notBefore" placeholder="RFC3339"></label>
          <label>Not After<input id="notAfter" placeholder="RFC3339"></label>
        </div><div class="grid3"><label>Pre-Check Mode<select id="stagePolicyPreCheck"><option value="first-match">first-match</option><option value="continue-all">continue-all</option></select></label><label>Rewrite Mode<select id="stagePolicyRewrite"><option value="first-match">first-match</option><option value="continue-all">continue-all</option></select></label><label>Route Mode<select id="stagePolicyRoute"><option value="first-match">first-match</option><option value="continue-all">continue-all</option></select></label></div><div class="grid2"><label>Mirror Mode<select id="stagePolicyMirror"><option value="continue-all">continue-all</option><option value="first-match">first-match</option></select></label><label>Post-Audit Mode<select id="stagePolicyPostAudit"><option value="continue-all">continue-all</option><option value="first-match">first-match</option></select></label></div>
        <div class="actions"><button class="warn" id="addRouteToModelRuleBtn">Add Route Rule</button><button class="warn" id="addFallbackRuleBtn">Add Fallback Rule</button><button class="warn" id="addDenyRuleBtn">Add Deny Rule</button><button class="warn" id="buildFallbackChainBtn">Build Fallback Chain Rule</button></div><div class="card" style="padding:12px 0 0;border:0;box-shadow:none;background:none"><h3>Templates</h3><div class="grid2"><label>Template Search<input id="templateSearchInput" placeholder="name or description"></label><label>Template Category<select id="templateCategoryFilter"><option value="">all</option><option value="routing">routing</option><option value="fallback">fallback</option><option value="security">security</option><option value="custom">custom</option></select></label></div><div class="grid3"><label>Scenario<select id="templateScenarioFilter"><option value="">all</option><option value="model-migration">model-migration</option><option value="traffic-split">traffic-split</option><option value="cost-control">cost-control</option><option value="shadow-release">shadow-release</option><option value="provider-guardrail">provider-guardrail</option></select></label><label>Maturity<select id="templateMaturityFilter"><option value="">all</option><option value="stable">stable</option><option value="beta">beta</option><option value="experimental">experimental</option></select></label><label>Tag<input id="templateTagFilter" placeholder="routing"></label></div><div class="actions"><button id="saveTemplateBtn">Save Current Rule As Template</button><button class="alt" id="applyTemplateFilterBtn">Filter Templates</button><button class="alt" id="exportTemplatesBtn">Export Templates</button></div><label>Template Import JSON<textarea id="templateImportBox" placeholder="{&quot;items&quot;:[...]}"></textarea></label><div class="actions"><button class="warn" id="importTemplatesBtn">Import Templates</button></div><div class="list" id="templateList"></div></div>
        <div class="rule-builder">
          <div class="grid4"><label>Rule ID<input id="ruleIdInput"></label><label>Stage<select id="ruleStageInput"><option value="pre-check">pre-check</option><option value="rewrite">rewrite</option><option value="route">route</option><option value="mirror">mirror</option><option value="post-audit">post-audit</option></select></label><label>Match Model<input id="ruleMatchModelInput" placeholder="gpt-5.5,gpt-5.4"></label><label>Priority<input id="rulePriorityInput" type="number" value="10"></label></div>
          <div class="grid3"><label>Match Path<input id="rulePathInput" placeholder="/v1/responses,/v1/chat/completions"></label><label>Match Provider<input id="ruleProviderInput" placeholder="openai,codex"></label><label>Match Header<input id="ruleHeaderInput" placeholder="X-Test:yes"></label></div>
          <div class="grid3"><label>Match Query<input id="ruleQueryInput" placeholder="mode:fast;tenant:a"></label><label>Body Contains<input id="ruleBodyContainsInput" placeholder="input.0.role:user;service_tier:priority"></label><label>Metadata Contains<input id="ruleMetadataInput" placeholder="client.tag:tenant-a"></label></div>
          <div class="section-title"><h3>Routing Actions</h3><span class="hint">Build direct routes, weighted splits, fallbacks, mirrors, and provider forcing.</span></div>
          <div class="split">
            <div class="route-board">
              <div class="grid3"><label>Route To Model<input id="ruleRouteToModelInput" placeholder="openai/gpt-5.4"></label><label>Force Provider Prefix<input id="ruleForceProviderInput" placeholder="openai"></label><label>Deny Provider<input id="ruleDenyProviderInput" placeholder="claude"></label></div>
              <div class="grid3"><label>Shard By<select id="ruleShardByInput"><option value="">default(api_key)</option><option value="api_key">api_key</option><option value="model">model</option><option value="path">path</option><option value="user_agent">user_agent</option><option value="client_ip">client_ip</option><option value="query">query</option><option value="header">header</option></select></label><label>Mirror Models<input id="ruleMirrorModelsInput" placeholder="openai/gpt-5.4-mini,codex/gpt-4.1-mini"></label><label>Reason Tag<input id="ruleReasonTagInput" placeholder="tenant-a"></label></div>
              <div class="grid3"><label>Fallback Models<input id="ruleFallbackModelInput" placeholder="openai/gpt-5.4-mini,openai/gpt-4.1-mini"></label><label>Weighted Route Model<input id="weightedRouteModelInput" placeholder="openai/gpt-5.4"></label><label>Weight<input id="weightedRouteWeightInput" type="number" min="1" value="1"></label></div><div class="grid4"><label>Route Pool Name<input id="routePoolNameInput" placeholder="primary-openai-pool"></label><label>Route Pool Mode<select id="routePoolModeInput"><option value="weighted">weighted</option><option value="shard">shard</option></select></label><label>Provider Affinity<input id="routePoolAffinityInput" placeholder="openai"></label><label>Failover Chain<input id="ruleFailoverChainInput" placeholder="openai/gpt-5.4,openai/gpt-4.1-mini"></label></div><div class="grid4"><label>Pool Member Provider<input id="weightedRouteProviderInput" placeholder="openai"></label><label>Pool Member Suffix<input id="weightedRouteSuffixInput" placeholder="gpt-5.4"></label><label>Member Priority<input id="weightedRoutePriorityInput" type="number" value="100"></label><label>Member Enabled<select id="weightedRouteEnabledInput"><option value="true">true</option><option value="false">false</option></select></label></div><div class="grid4"><label>Member Status<select id="weightedRouteStatusInput"><option value="active">active</option><option value="degraded">degraded</option><option value="drain">drain</option><option value="offline">offline</option></select></label><label>Member Reason<input id="weightedRouteReasonInput" placeholder="manual-drain"></label><label>Health Score<input id="weightedRouteHealthInput" type="number" min="0" max="100" value="100"></label><label>Traffic Cap<input id="weightedRouteTrafficCapInput" type="number" min="0" max="100" value="100"></label></div><div class="mini-card"><h4>Failover Hops</h4><div class="grid4"><label>Hop Model<input id="failoverHopModelInput" placeholder="openai/gpt-5.4"></label><label>Hop Provider<input id="failoverHopProviderInput" placeholder="openai"></label><label>Hop Suffix<input id="failoverHopSuffixInput" placeholder="gpt-5.4"></label><label>On Decision<select id="failoverHopDecisionInput"><option value="reject">reject</option><option value="rewrite">rewrite</option><option value="pass">pass</option></select></label></div><div class="grid3"><label>Hop Reason<input id="failoverHopReasonInput" placeholder="quota-exceeded"></label><label>Hop Enabled<select id="failoverHopEnabledInput"><option value="true">true</option><option value="false">false</option></select></label><label>Advanced Hop JSON<textarea id="ruleFailoverHopsInput" placeholder="[{&quot;provider&quot;:&quot;openai&quot;,&quot;suffix&quot;:&quot;gpt-5.4&quot;,&quot;reason&quot;:&quot;quota-exceeded&quot;,&quot;on_decision&quot;:&quot;reject&quot;}]"></textarea></label></div><div class="actions"><button class="alt" id="addFailoverHopBtn">Add Failover Hop</button><button class="alt" id="clearFailoverHopsBtn">Clear Hops</button></div><div class="tablelist" id="failoverHopsList"></div></div><div class="hint">Pool members can be full models or provider+suffix pairs.</div>
              <div class="actions"><button class="alt" id="addWeightedRouteBtn">Add Weighted Route</button><button class="alt" id="clearWeightedRoutesBtn">Clear Weighted Routes</button><button class="alt" id="sortWeightedRoutesBtn">Sort Weighted Routes</button></div>
              <div class="tablelist" id="weightedRoutesList"></div><div class="mini-card"><h4>Pool Health</h4><div id="poolHealthOut" class="list"></div></div><div class="mini-card"><h4>Route Graph</h4><div id="routeGraphOut" class="list"></div></div>
              <div class="card soft" style="padding:12px 0 0;border:0;box-shadow:none"><h3>Fallback Chain Preview</h3><div class="actions"><button class="alt" id="addFallbackHopBtn">Add Hop</button><button class="alt" id="removeFallbackHopBtn">Remove Last Hop</button><button class="alt" id="sortFallbackHopsBtn">Sort Hops</button></div><pre id="fallbackPreview">[]</pre></div>
            </div>
            <div class="group-board">
              <div class="mini-card"><h4>Condition Groups</h4><div class="grid2"><label>Use Any-Of<select id="ruleUseAnyOf"><option value="false">false</option><option value="true">true</option></select></label><label>Use All-Of<select id="ruleUseAllOf"><option value="false">false</option><option value="true">true</option></select></label></div><div class="actions"><button class="alt" id="buildAnyOfBtn">Append Any-Of Group</button><button class="alt" id="buildAllOfBtn">Append All-Of Group</button><button class="alt" id="syncCurrentRuleToGroupsBtn">Sync Current Match</button></div><div class="actions"><button class="alt" id="loadFirstAnyOfBtn">Load First Any-Of</button><button class="alt" id="loadFirstAllOfBtn">Load First All-Of</button><button class="alt" id="popConditionGroupBtn">Remove Last Group</button><button class="alt" id="clearConditionGroupsBtn">Clear Groups</button></div><div id="conditionGroupsPreview" class="list"></div></div>
              <label>Any-Of JSON<textarea id="ruleAnyOfInput" placeholder="[{&quot;models&quot;:[&quot;gpt-5.5&quot;]},{&quot;paths&quot;:[&quot;/v1/responses&quot;]}]"></textarea></label>
              <label>All-Of JSON<textarea id="ruleAllOfInput" placeholder="[{&quot;query&quot;:{&quot;mode&quot;:&quot;strict&quot;}},{&quot;body_contains&quot;:{&quot;service_tier&quot;:&quot;priority&quot;}}]"></textarea></label>
            </div>
          </div>
          <div class="actions"><button class="warn" id="pushRuleFromFormBtn">Append Rule From Form</button><button id="saveRuleBtn">Save Rule</button><button id="cloneRuleBtn">Clone Rule</button><button class="danger" id="deleteRuleBtn">Delete Rule</button></div><div class="actions"><button class="alt" id="moveRuleUpBtn">Move Rule Up</button><button class="alt" id="moveRuleDownBtn">Move Rule Down</button></div><div class="list" id="rulesList"></div><label>Rules JSON<textarea id="rulesBox"></textarea></label>
        </div>
        <div class="actions">
          <button id="savePolicyBtn">Save Selected Policy</button>
          <button class="alt" id="exportPolicyBundleBtn">Export Policy Bundle</button>
          <button class="danger" id="deletePolicyBtn">Delete Selected Policy</button>
          <button class="danger" id="resetUsageBtn">Reset Selected Usage</button>
        </div>
      </article>
      <article class="card" style="position:sticky;top:16px;z-index:1">
        <h2>Current Node</h2>
        <div id="currentNodeDetail" class="list"></div>
      </article><article class="card">
        <h2>Audit Detail</h2>
        <div class="stats" id="auditDetailStats"></div>
        <div class="stats" id="auditTimelineStats"></div>
        <pre id="auditDetailOut">{}</pre>
      </article>
      <article class="card">
        <h2>Dry Run</h2>
        <div class="grid3">
          <label>API Key<input id="dryKey" placeholder="top-level api key"></label>
          <label>Model<input id="dryModel" value="gpt-5.5"></label>
          <label>Source Format<select id="dryFormat"><option value="openai-response">openai-response</option><option value="openai">openai</option><option value="claude">claude</option></select></label>
        </div>
        <div class="grid2">
          <label>Request Path<input id="dryPath" value="/v1/responses"></label>
          <label>Stream<select id="dryStream"><option value="false">false</option><option value="true">true</option></select></label>
        </div>
        <label>Request Body<textarea id="dryBody">{"model":"gpt-5.5","messages":[{"role":"user","content":"hello"}]}</textarea></label>
        <div class="actions"><button class="warn" id="dryRunBtn">Run Dry-Run</button><button class="alt" id="copyDryRunHintsBtn">Copy Hints</button><button class="alt" id="clearDryRunHintsBtn">Clear Hints</button></div>
        <div class="stats" id="dryRunHintStats"></div>
        <div class="stats" id="dryRunDecisionStats"></div>
        <div class="stats" id="dryRunStageTrace"></div>
        <pre id="dryRunOut">{}</pre>
      </article>
      <article class="card">
        <h2>Advanced JSON</h2>
        <p class="hint">Use this when you want to edit the full policy document directly.</p>
        <div class="actions"><button id="saveAllBtn">Save Whole Policy Set</button><button class="alt" id="importPolicyBundleBtn">Import Bundle</button></div>
        <div class="grid2"><label>Bundle Mode<select id="policyBundleMode"><option value="merge">merge</option><option value="replace">replace</option></select></label><label>Bundle Name<input id="policyBundleName" placeholder="gateway-policy-bundle"></label></div>
        <label>Policies JSON<textarea id="policiesBox" class="big"></textarea></label>
      </article>
    </main>
  </section>
</div>
<script nonce="gateway-ui">
const api = {
  keys: '/v0/management/plugins/gateway/keys',
  health: '/v0/management/plugins/gateway/health',
  policies: '/v0/management/plugins/gateway/policies',
  exportPolicies: '/v0/management/plugins/gateway/policies/export',
  importPolicies: '/v0/management/plugins/gateway/policies/import',
  clonePolicy: '/v0/management/plugins/gateway/policies/clone',
  usage: '/v0/management/plugins/gateway/usage',
  audit: '/v0/management/plugins/gateway/audit',
  auditDetail: '/v0/management/plugins/gateway/audit/detail',
  auditSummary: '/v0/management/plugins/gateway/audit/summary',
  templates: '/v0/management/plugins/gateway/templates',
  exportTemplates: '/v0/management/plugins/gateway/templates/export',
  importTemplates: '/v0/management/plugins/gateway/templates/import',
  addPolicy: '/v0/management/plugins/gateway/policies/add',
  addRule: '/v0/management/plugins/gateway/rules/add',
  rule: '/v0/management/plugins/gateway/rules',
  resetUsage: '/v0/management/plugins/gateway/usage/reset',
  routeMemberOp: '/v0/management/plugins/gateway/route-members/op',
  routeMemberPreview: '/v0/management/plugins/gateway/route-members/preview',
  dryRun: '/v0/management/plugins/gateway/dry-run'
};
const gatewayTokenParam = new URLSearchParams(window.location.search).get('gateway_token') || '';
const managementKeyStorage = 'gateway-management-key';
let managementAuthBlocked = false;
function apiURL(url){
  if(!gatewayTokenParam) return url;
  return url + (url.includes('?') ? '&' : '?') + 'gateway_token=' + encodeURIComponent(gatewayTokenParam);
}
function loadManagementKey(){ try { return window.sessionStorage.getItem(managementKeyStorage) || ''; } catch { return ''; } }
function storeManagementKey(value){ try { if(value) window.sessionStorage.setItem(managementKeyStorage, value); else window.sessionStorage.removeItem(managementKeyStorage); } catch {} }
function currentManagementKey(){ return (el('managementKeyInput')?.value || loadManagementKey() || '').trim(); }
function setManagementStatus(message, kind='ok'){
  const node = el('managementAuthStatus');
  if(!node) return;
  node.value = message;
  node.style.color = kind === 'bad' ? '#b42318' : '#0f766e';
}
function managementInit(init){
  const next = Object.assign({}, init || {});
  const headers = new Headers((init && init.headers) || {});
  const key = currentManagementKey();
  if(key){
    headers.set('Authorization', 'Bearer ' + key);
    headers.set('X-Management-Key', key);
  }
  next.headers = headers;
  next.credentials = 'same-origin';
  return next;
}
function ensureManagementAccess(){
  if(!currentManagementKey()){
    setManagementStatus('Enter management key first', 'bad');
    throw new Error('Enter the CPAMC management key before connecting.');
  }
  if(managementAuthBlocked){
    setManagementStatus('Auth failed; save key or reconnect', 'bad');
    throw new Error('Management authentication is blocked after a failed attempt. Save the key or reconnect before retrying.');
  }
}
const rawInnerHTMLDescriptor = Object.getOwnPropertyDescriptor(Element.prototype, 'innerHTML');
function sanitizeHTML(raw){
  const template = document.createElement('template');
  rawInnerHTMLDescriptor.set.call(template, String(raw ?? ''));
  template.content.querySelectorAll('script,iframe,object,embed,link,meta,style,svg,math,form').forEach(node => node.remove());
  template.content.querySelectorAll('*').forEach(node => {
    Array.from(node.attributes).forEach(attribute => {
      const name = attribute.name.toLowerCase();
      const value = String(attribute.value || '').trim().toLowerCase();
      if(name.startsWith('on') || name === 'style' || ((name === 'href' || name === 'src' || name === 'srcdoc' || name === 'xlink:href') && value.startsWith('javascript:'))){
        node.removeAttribute(attribute.name);
      }
    });
  });
  return rawInnerHTMLDescriptor.get.call(template);
}
Object.defineProperty(Element.prototype, 'innerHTML', {
  get(){ return rawInnerHTMLDescriptor.get.call(this); },
  set(value){ rawInnerHTMLDescriptor.set.call(this, sanitizeHTML(value)); }
});
const el = id => document.getElementById(id);
function clearNode(node){ if(node) node.replaceChildren(); }
function textNode(tag, className, text){
  const node = document.createElement(tag);
  if(className) node.className = className;
  node.textContent = String(text ?? '');
  return node;
}
function appendText(parent, tag, className, text){
  const node = textNode(tag, className, text);
  parent.appendChild(node);
  return node;
}
function appendButton(parent, className, text, onClick){
  const button = document.createElement('button');
  if(className) button.className = className;
  button.type = 'button';
  button.textContent = text;
  if(onClick) button.addEventListener('click', onClick);
  parent.appendChild(button);
  return button;
}
function emptyNode(className, text){ return textNode('div', className, text); }
function cardTitle(card){ return (card?.querySelector('h2')?.textContent || '').trim(); }
function buildTabbedLayout(){
  const layout = document.querySelector('.layout');
  if(!layout || layout.classList.contains('tab-shell')) return;
  const stacks = Array.from(layout.children).filter(node => node.classList && node.classList.contains('stack'));
  const cards = [];
  stacks.forEach(stack => {
    Array.from(stack.children).forEach(child => {
      if(child.matches && child.matches('article.card')) cards.push(child);
    });
  });
  if(!cards.length) return;
  const byTitle = new Map();
  cards.forEach(card => {
    const title = cardTitle(card);
    if(title && !byTitle.has(title)) byTitle.set(title, card);
  });
  const groups = [
    { id: 'dashboard', label: 'Dashboard', hint: 'keys, health, activity', columns: '2', titles: ['Configured Keys', 'Diagnostics', 'Activity', 'Guide'] },
    { id: 'traffic', label: 'Traffic', hint: 'usage, topology, node', columns: '2', titles: ['Usage Snapshot', 'Current Node'] },
    { id: 'policies', label: 'Policies', hint: 'limits, rules, JSON', columns: '', titles: ['Selected Policy', 'Advanced JSON'] },
    { id: 'audit', label: 'Audit', hint: 'events and evidence', columns: 'audit', titles: ['Audit Summary', 'Audit Log', 'Audit Detail'] },
    { id: 'dryrun', label: 'Dry Run', hint: 'simulate routing', columns: '', titles: ['Dry Run'] }
  ];
  const used = new Set();
  const nav = document.createElement('nav');
  nav.className = 'tab-nav';
  nav.setAttribute('aria-label', 'Gateway management sections');
  const content = document.createElement('div');
  content.className = 'tab-content';
  function activateTab(id){
    nav.querySelectorAll('.tab-button').forEach(button => {
      const active = button.dataset.tabTarget === id;
      button.classList.toggle('active', active);
      button.setAttribute('aria-selected', active ? 'true' : 'false');
    });
    content.querySelectorAll('.tab-pane').forEach(pane => pane.classList.toggle('active', pane.id === 'gateway-tab-' + id));
    try { window.sessionStorage.setItem('gateway-active-tab', id); } catch {}
  }
  groups.forEach(group => {
    const pane = document.createElement('section');
    pane.className = 'tab-pane';
    pane.id = 'gateway-tab-' + group.id;
    pane.setAttribute('role', 'tabpanel');
    if(group.columns) pane.dataset.columns = group.columns;
    group.titles.forEach(title => {
      const card = byTitle.get(title);
      if(!card) return;
      if(title === 'Current Node') card.removeAttribute('style');
      card.dataset.gatewaySection = group.id;
      pane.appendChild(card);
      used.add(card);
    });
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'tab-button';
    button.dataset.tabTarget = group.id;
    button.setAttribute('role', 'tab');
    appendText(button, 'b', '', group.label);
    appendText(button, 'span', '', group.hint);
    button.addEventListener('click', () => activateTab(group.id));
    nav.appendChild(button);
    content.appendChild(pane);
  });
  const leftovers = cards.filter(card => !used.has(card));
  if(leftovers.length){
    const pane = document.createElement('section');
    pane.className = 'tab-pane';
    pane.id = 'gateway-tab-more';
    pane.setAttribute('role', 'tabpanel');
    leftovers.forEach(card => pane.appendChild(card));
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'tab-button';
    button.dataset.tabTarget = 'more';
    button.setAttribute('role', 'tab');
    appendText(button, 'b', '', 'More');
    appendText(button, 'span', '', 'extra panels');
    button.addEventListener('click', () => activateTab('more'));
    nav.appendChild(button);
    content.appendChild(pane);
  }
  layout.className = 'layout tab-shell';
  layout.replaceChildren(nav, content);
  let initial = 'dashboard';
  try { initial = window.sessionStorage.getItem('gateway-active-tab') || initial; } catch {}
  if(!content.querySelector('#gateway-tab-' + initial)) initial = 'dashboard';
  activateTab(initial);
}
function appendStat(parent, label, value){
  const node = document.createElement('div');
  node.className = 'stat';
  appendText(node, 'span', '', label);
  appendText(node, 'b', '', value);
  parent.appendChild(node);
  return node;
}
function appendPreviewChip(parent, type, label, active){
  const chip = textNode('span', 'chip', label);
  chip.dataset.previewFocus = type;
  chip.dataset.previewLabel = label;
  if(active){
    chip.style.outline = '2px solid #a16207';
    chip.style.outlineOffset = '2px';
  }
  parent.appendChild(chip);
  return chip;
}
function appendNodeDetail(parent, title, lines){
  const detail = document.createElement('div');
  detail.className = 'node-detail';
  appendText(detail, 'b', '', title);
  (lines || []).forEach(line => appendText(detail, 'span', '', line));
  parent.appendChild(detail);
  return detail;
}
function appendSubnodeList(parent, type, items, options){
  const list = Array.isArray(items) ? items : [];
  const emptyText = options?.emptyText || '';
  if(!list.length){
    if(emptyText) appendText(parent, 'div', 'hint', emptyText);
    return null;
  }
  const out = document.createElement('div');
  out.className = 'subnode-list';
  list.forEach((item, index) => {
    const label = previewChipLabel(type, item);
    const node = document.createElement('div');
    node.className = 'subnode ' + (options?.kind || '') + (previewMatches(type, label) ? ' focused' : '');
    node.dataset.previewFocus = type;
    node.dataset.previewLabel = label;
    const body = document.createElement('div');
    appendText(body, 'strong', '', label);
    const summary = options?.summary ? options.summary(item, index) : '';
    if(summary) appendText(body, 'span', 'hint', summary);
    node.appendChild(body);
    appendText(node, 'div', 'hint', options?.meta ? options.meta(item, index) : '');
    out.appendChild(node);
  });
  parent.appendChild(out);
  return out;
}
const safeApplyInput = el('poolSafeApplyInput');
if(safeApplyInput){ safeApplyInput.value = 'true'; safeApplyInput.disabled = true; }
const initialURLState = readURLState();
let state = { keys: [], usage: [], health: {}, audit: [], auditSummary: { total_by_decision: {}, total_by_reason: {}, total_by_rule: {}, total_by_policy: {}, total_by_model: {} }, templates: [], policies: { key_policies: [], default_policy: {} }, selectedKeyId: initialURLState.keyId || '', selectedRuleId: initialURLState.ruleId || '', focusedPreviewLabel: initialURLState.previewLabel || '', focusedPreviewType: initialURLState.previewType || '', latestPoolPreviewToken: '', weightedRoutes: [], failoverHops: [], memberHitCounts: {}, ruleHitCounts: {}, stageHitCounts: {}, memberHitCountsLast5m: {}, ruleHitCountsLast5m: {}, stageHitCountsLast5m: {}, memberHitCountsLastHour: {}, ruleHitCountsLastHour: {}, stageHitCountsLastHour: {}, memberHitCountsLast24h: {}, ruleHitCountsLast24h: {}, stageHitCountsLast24h: {}, statsWindow: initialURLState.statsWindow || '1h', collapsedTopology: Object.keys(initialURLState.collapsedTopology || {}).length ? initialURLState.collapsedTopology : loadCollapsedTopology() };
function log(msg, kind='ok'){ el('logBox').textContent = '[' + new Date().toLocaleTimeString() + '] ' + msg; el('logBox').style.color = kind === 'bad' ? '#ffb4ab' : '#9ff3e8'; }
async function readJSON(url, init){
  if(url.startsWith('/v0/management/')) ensureManagementAccess();
  const res = await fetch(apiURL(url), managementInit(init));
  const text = await res.text();
  let body = {};
  if(text){
    try { body = JSON.parse(text); } catch { body = { message: text }; }
  }
  if(res.status === 401 || res.status === 403){
    managementAuthBlocked = true;
    setManagementStatus('Auth failed: check management key', 'bad');
    throw new Error(body.error || body.message || 'Management authentication failed.');
  }
  if(!res.ok) throw new Error(body.error || body.message || JSON.stringify(body));
  return body;
}
function selectedPolicy(){ return (state.policies.key_policies || []).find(item => item.key_id === state.selectedKeyId) || null; }
function renderKeys(){
  const root = el('keysList'); clearNode(root);
  const items = state.keys || [];
  if(!items.length){ root.appendChild(emptyNode('item', 'No key-specific policies yet.')); return; }
  items.forEach(key => {
    const node = document.createElement('div'); node.className='item' + (state.selectedKeyId === key.key_id ? ' active' : '');
    appendText(node, 'strong', '', key.display_name || key.key_id);
    appendText(node, 'span', 'pill', key.masked_key || '');
    appendText(node, 'div', 'hint', 'key_id: ' + key.key_id + ' | enabled: ' + key.enabled);
    node.addEventListener('click', () => { state.selectedKeyId = key.key_id; state.selectedRuleId = ''; syncURLState(); hydrateSelectedPolicy(); renderKeys(); renderRules(); });
    root.appendChild(node);
  });
}
function topologyCollapsed(key){ return Boolean(state.collapsedTopology && state.collapsedTopology[key]); }
function loadCollapsedTopology(){ try { const raw = window.localStorage.getItem('gateway-topology-collapse'); if(!raw) return {}; const parsed = JSON.parse(raw); return parsed && typeof parsed === 'object' ? parsed : {}; } catch { return {}; } }
function readURLState(){ try { const params = new URLSearchParams(window.location.search); const collapsed = params.get('collapsed'); return { statsWindow: params.get('window') || '1h', keyId: params.get('key') || '', ruleId: params.get('rule') || '', previewLabel: params.get('preview') || '', previewType: params.get('previewType') || '', collapsedTopology: collapsed ? Object.fromEntries(collapsed.split(',').filter(Boolean).map(key => [key, true])) : {} }; } catch { return { statsWindow: '1h', keyId: '', ruleId: '', previewLabel: '', previewType: '', collapsedTopology: {} }; } }
function syncURLState(){ try { const params = new URLSearchParams(window.location.search); if(state.statsWindow) params.set('window', state.statsWindow); else params.delete('window'); if(state.selectedKeyId) params.set('key', state.selectedKeyId); else params.delete('key'); if(state.selectedRuleId) params.set('rule', state.selectedRuleId); else params.delete('rule'); if(state.focusedPreviewLabel) params.set('preview', state.focusedPreviewLabel); else params.delete('preview'); if(state.focusedPreviewType) params.set('previewType', state.focusedPreviewType); else params.delete('previewType'); const collapsedKeys = Object.entries(state.collapsedTopology || {}).filter(([, value]) => value).map(([key]) => key); if(collapsedKeys.length) params.set('collapsed', collapsedKeys.join(',')); else params.delete('collapsed'); const query = params.toString(); const next = window.location.pathname + (query ? '?' + query : ''); window.history.replaceState(null, '', next); } catch {} }
function saveCollapsedTopology(){ try { window.localStorage.setItem('gateway-topology-collapse', JSON.stringify(state.collapsedTopology || {})); } catch {} }
function toggleTopology(key){ if(!state.collapsedTopology) state.collapsedTopology = {}; state.collapsedTopology[key] = !state.collapsedTopology[key]; saveCollapsedTopology(); renderGlobalTopology(); }
function previewTargetLabel(item){ return item?.model || ((item?.provider || '-') + '/' + (item?.suffix || '-')); }
function previewChipLabel(type, item){
  const target = previewTargetLabel(item);
  if(type === 'pool-member') return target + ' [' + (item?.status || 'active') + ']';
  if(type === 'failover-hop') return target + ' [' + (item?.reason || 'fallback') + ']';
  if(type === 'mirror-target') return target + ' [mirror]';
  return target;
}
function previewMatches(type, label){ return state.focusedPreviewType === (type || '') && state.focusedPreviewLabel === (label || ''); }
function setPreviewFocus(type, label){ state.focusedPreviewType = type || ''; state.focusedPreviewLabel = label || ''; syncURLState(); }
function clearPreviewFocus(){ state.focusedPreviewType = ''; state.focusedPreviewLabel = ''; syncURLState(); }
function bindPreviewFocus(scope, policy, rule){
  if(!scope) return;
  scope.querySelectorAll('[data-preview-focus]').forEach(node => {
    node.addEventListener('click', (event) => {
      event.stopPropagation();
      const type = node.getAttribute('data-preview-focus') || '';
      const label = node.getAttribute('data-preview-label') || '';
      setPreviewFocus(type, label);
      if(policy?.key_id) state.selectedKeyId = policy.key_id;
      if(rule?.id) state.selectedRuleId = rule.id;
      if(policy) hydrateSelectedPolicy();
      if(rule) hydrateRule(rule);
      renderKeys();
      renderRules();
      renderRouteGraph();
      renderCurrentNodeDetail();
      renderGlobalTopology();
      log('Focused topology target: ' + label + '.');
    });
  });
}
function statTone(value){ if(value >= 20) return 'warn'; if(value >= 5) return 'pill'; return 'hint'; }

function currentWindowHitMaps(){
  switch(state.statsWindow){
    case '5m': return { member: state.memberHitCountsLast5m || {}, rule: state.ruleHitCountsLast5m || {}, stage: state.stageHitCountsLast5m || {}, label: '5m' };
    case '24h': return { member: state.memberHitCountsLast24h || {}, rule: state.ruleHitCountsLast24h || {}, stage: state.stageHitCountsLast24h || {}, label: '24h' };
    default: return { member: state.memberHitCountsLastHour || {}, rule: state.ruleHitCountsLastHour || {}, stage: state.stageHitCountsLastHour || {}, label: '1h' };
  }
}

function currentDryRunHintPayload(){
  try { return JSON.parse(el('dryBody').value || '{}'); } catch { return null; }
}

function buildDryRunPayloadFromRule(policy, rule){
  if(!rule) return null;
  return {
    model: rule?.match?.models?.[0] || el('dryModel').value,
    messages: [{ role: 'user', content: 'gateway dry-run from current rule' }],
    route_hint: rule?.actions?.route_to_model || '',
    provider_hint: (rule?.match?.providers || [])[0] || '',
    headers_hint: rule?.match?.headers || {},
    query_hint: rule?.match?.query || {},
    body_contains_hint: rule?.match?.body_contains || {},
    metadata_hint: rule?.match?.metadata_contains || {},
    failover_chain_hint: rule?.actions?.failover_chain || [],
    mirror_models_hint: rule?.actions?.mirror_models || []
  };
}

function applyDryRunPayloadFromRule(policy, rule){
  if(!rule) return;
  if(policy?.match_api_key) el('dryKey').value = policy.match_api_key;
  else if(policy?.key_id) el('dryKey').value = '';
  if(rule?.match?.models?.[0]) el('dryModel').value = rule.match.models[0];
  if(rule?.match?.paths?.[0]) el('dryPath').value = rule.match.paths[0];
  const payload = buildDryRunPayloadFromRule(policy, rule);
  if(!payload) return;
  el('dryBody').value = JSON.stringify(payload, null, 2);
  renderDryRunHints(payload);
}

function renderDryRunHints(payload){
  const root = el('dryRunHintStats');
  if(!root) return;
  clearNode(root);
  if(!payload || typeof payload !== 'object'){ appendStat(root, 'No dry-run hints yet.', '0'); return; }
  [
    ['provider_hint', payload.provider_hint],
    ['route_hint', payload.route_hint],
    ['query_hint', payload.query_hint && Object.keys(payload.query_hint).length ? JSON.stringify(payload.query_hint) : ''],
    ['headers_hint', payload.headers_hint && Object.keys(payload.headers_hint).length ? JSON.stringify(payload.headers_hint) : ''],
    ['body_hint', payload.body_contains_hint && Object.keys(payload.body_contains_hint).length ? JSON.stringify(payload.body_contains_hint) : ''],
    ['metadata_hint', payload.metadata_hint && Object.keys(payload.metadata_hint).length ? JSON.stringify(payload.metadata_hint) : ''],
    ['failover_hint', payload.failover_chain_hint && payload.failover_chain_hint.length ? payload.failover_chain_hint.join(', ') : ''],
    ['mirror_hint', payload.mirror_models_hint && payload.mirror_models_hint.length ? payload.mirror_models_hint.join(', ') : '']
  ].forEach(([label, value]) => {
    if(!value) return;
    appendStat(root, label, value);
  });
  if(!root.childElementCount){ appendStat(root, 'No dry-run hints yet.', '0'); }
}

function renderDryRunDecision(result){
  const root = el('dryRunDecisionStats');
  if(!root) return;
  clearNode(root);
  if(!result || typeof result !== 'object'){
    appendStat(root, 'No dry-run decision yet.', '0');
    return;
  }
  [
    ['decision', result.decision || 'pass'],
    ['reason', result.reason || '-'],
    ['final_model', result.final_model || '-'],
    ['rule_id', result.rule_id || '-']
  ].forEach(([label, value]) => {
    appendStat(root, label, value);
  });
}

function renderDryRunStageTrace(items){
  const root = el('dryRunStageTrace');
  if(!root) return;
  clearNode(root);
  const rows = Array.isArray(items) ? items : [];
  if(!rows.length){
    appendStat(root, 'No stage trace yet.', '0');
    return;
  }
  rows.forEach((item, index) => {
    const node = document.createElement('div');
    node.className = 'item';
    const failoverReasons = Array.isArray(item.failover_reasons) && item.failover_reasons.length ? item.failover_reasons.join(', ') : '-';
    const mirrors = Array.isArray(item.mirror_models) && item.mirror_models.length ? item.mirror_models.join(', ') : '-';
    const failoverChain = Array.isArray(item.failover_chain) && item.failover_chain.length ? item.failover_chain.join(' -> ') : '-';
    const matchedCount = item.matched_count ?? ((item.matched_rules || []).length || 0);
    appendText(node, 'strong', '', (index + 1) + '. ' + (item.stage || '-'));
    appendText(node, 'div', 'hint', 'mode ' + (item.mode || '-') + ' | matched ' + matchedCount + ' | decision ' + (item.decision || 'pass'));
    appendText(node, 'div', 'hint', 'rules ' + (((item.matched_rules || []).length ? item.matched_rules.join(', ') : '-')));
    appendText(node, 'div', 'hint', 'final model ' + (item.final_model || '-') + ' | route target ' + (item.route_target || '-'));
    appendText(node, 'div', 'hint', 'route pool ' + (item.route_pool || '-') + ' | fallback target ' + (item.fallback_target || '-'));
    appendText(node, 'div', 'hint', 'mirrors ' + mirrors);
    appendText(node, 'div', 'hint', 'failover chain ' + failoverChain);
    appendText(node, 'div', 'hint', 'failover reasons ' + failoverReasons);
    appendText(node, 'div', 'hint', 'reason ' + (item.reason || '-'));
    root.appendChild(node);
  });
}

function renderCurrentNodeDetail(){
  const root = el('currentNodeDetail');
  if(!root) return;
  clearNode(root);
  const policy = selectedPolicy();
  const rule = (policy?.rules || []).find(item => item.id === state.selectedRuleId);
  if(!policy){ root.appendChild(emptyNode('item', 'No policy selected.')); return; }
  const windowHits = currentWindowHitMaps();
  const policyCard = document.createElement('div');
  policyCard.className = 'item';
  appendText(policyCard, 'strong', '', 'Policy');
  appendText(policyCard, 'div', 'hint', policy.display_name || policy.key_id);
  appendText(policyCard, 'div', 'hint', 'enabled ' + (policy.enabled !== false) + ' | rules ' + ((policy.rules || []).length) + ' | window ' + windowHits.label);
  root.appendChild(policyCard);
  const governanceCard = document.createElement('div');
  governanceCard.className = 'item';
  appendText(governanceCard, 'strong', '', 'Governance');
  appendText(governanceCard, 'div', 'hint', 'requests/day ' + (policy.limits?.requests_per_day || 0) + ' | requests/min ' + (policy.limits?.requests_per_minute || 0) + ' | inflight ' + (policy.limits?.max_inflight || 0));
  appendText(governanceCard, 'div', 'hint', 'pre-check ' + (policy.stage_policy?.['pre-check']?.mode || 'first-match') + ' | route ' + (policy.stage_policy?.route?.mode || 'first-match') + ' | mirror ' + (policy.stage_policy?.mirror?.mode || 'continue-all'));
  root.appendChild(governanceCard);
  const quickCard = document.createElement('div');
  quickCard.className = 'item';
  appendText(quickCard, 'strong', '', 'Quick Actions');
  const quickRow1 = document.createElement('div');
  quickRow1.className = 'actions';
  const focusPolicyBtn = appendButton(quickRow1, 'alt', 'Focus Policy', () => { renderKeys(); renderRules(); renderRouteGraph(); log('Focused current policy in topology.'); });
  focusPolicyBtn.id = 'focusPolicyBtn';
  const focusRuleBtn = appendButton(quickRow1, 'alt', 'Focus Rule', () => { if(rule){ hydrateRule(rule); renderRules(); renderRouteGraph(); log('Focused current rule.'); } });
  focusRuleBtn.id = 'focusRuleBtn';
  const openDryRunBtn = appendButton(quickRow1, 'alt', 'Open Dry-Run', () => {
    el('dryKey').value = policy.match_api_key || '';
    if(rule?.match?.models?.[0]) el('dryModel').value = rule.match.models[0];
    if(rule?.match?.paths?.[0]) el('dryPath').value = rule.match.paths[0];
    applyDryRunPayloadFromRule(policy, rule);
    log('Prepared dry-run form from current node.');
  });
  openDryRunBtn.id = 'openDryRunBtn';
  appendButton(quickRow1, 'alt', 'Reset Dry-Run').id = 'resetDryRunBtn';
  quickCard.appendChild(quickRow1);
  const quickRow2 = document.createElement('div');
  quickRow2.className = 'actions';
  appendButton(quickRow2, 'alt', 'Copy Request').id = 'copyDryRunRequestBtn';
  appendButton(quickRow2, 'alt', 'Copy Summary').id = 'copyRuleSummaryBtn';
  appendButton(quickRow2, 'alt', 'Copy Rule ID').id = 'copyRuleIdBtn';
  appendButton(quickRow2, 'alt', 'Copy Pool').id = 'copyRoutePoolBtn';
  appendButton(quickRow2, 'alt', 'Copy Failover').id = 'copyFailoverChainBtn';
  quickCard.appendChild(quickRow2);
  root.appendChild(quickCard);
  if(rule){
    const ruleCard = document.createElement('div');
    ruleCard.className = 'item';
    appendText(ruleCard, 'strong', '', 'Rule');
    appendText(ruleCard, 'div', 'hint', (rule.id || '-') + ' | stage ' + (rule.stage || 'pre-check'));
    appendText(ruleCard, 'div', 'hint', 'priority ' + (rule.priority || 0) + ' | total hits ' + (state.ruleHitCounts[rule.id || ''] || 0) + ' | current window ' + (windowHits.rule[rule.id || ''] || 0));
    root.appendChild(ruleCard);
    const actionCard = document.createElement('div');
    actionCard.className = 'item';
    appendText(actionCard, 'strong', '', 'Route');
    appendText(actionCard, 'div', 'hint', 'pool ' + (rule.actions?.route_pool?.name || '-') + ' | route_to ' + (rule.actions?.route_to_model || '-'));
    appendText(actionCard, 'div', 'hint', 'failover hops ' + ((rule.actions?.failover_hops || []).length) + ' | mirrors ' + ((rule.actions?.mirror_models || []).length));
    root.appendChild(actionCard);
    const poolMembers = rule.actions?.route_pool?.members || rule.actions?.weighted_routes || [];
    if(poolMembers.length){
      const poolCard = document.createElement('div');
      poolCard.className = 'item';
      appendText(poolCard, 'strong', '', 'Pool Members');
      const chips = document.createElement('div');
      chips.className = 'chips';
      poolMembers.slice(0, 3).forEach(member => {
        const label = previewChipLabel('pool-member', member);
        appendPreviewChip(chips, 'pool-member', label, previewMatches('pool-member', label));
      });
      if(poolMembers.length > 3) appendText(chips, 'span', 'chip', '...');
      poolCard.appendChild(chips);
      appendText(poolCard, 'div', 'hint', 'Click a member chip to refocus the topology route branch for the selected rule.');
      poolCard.title = 'Use the route branch in topology to inspect pool members for the selected rule.';
      bindPreviewFocus(poolCard, policy, rule);
      root.appendChild(poolCard);
    }
    if((rule.actions?.failover_hops || []).length){
      const hopCard = document.createElement('div');
      hopCard.className = 'item';
      appendText(hopCard, 'strong', '', 'Failover Hops');
      const chips = document.createElement('div');
      chips.className = 'chips';
      (rule.actions.failover_hops || []).slice(0, 3).forEach(hop => {
        const label = previewChipLabel('failover-hop', hop);
        appendPreviewChip(chips, 'failover-hop', label, previewMatches('failover-hop', label));
      });
      if((rule.actions.failover_hops || []).length > 3) appendText(chips, 'span', 'chip', '...');
      hopCard.appendChild(chips);
      appendText(hopCard, 'div', 'hint', 'Click a hop chip to localize the exact failover target in topology order.');
      hopCard.title = 'Use the failover branch in topology to inspect hop order for the selected rule.';
      bindPreviewFocus(hopCard, policy, rule);
      root.appendChild(hopCard);
    }
    if((rule.actions?.mirror_models || []).length){
      const mirrorCard = document.createElement('div');
      mirrorCard.className = 'item';
      appendText(mirrorCard, 'strong', '', 'Mirror Targets');
      const chips = document.createElement('div');
      chips.className = 'chips';
      (rule.actions.mirror_models || []).slice(0, 4).forEach(model => {
        const label = previewChipLabel('mirror-target', { model });
        appendPreviewChip(chips, 'mirror-target', label, previewMatches('mirror-target', label));
      });
      if((rule.actions.mirror_models || []).length > 4) appendText(chips, 'span', 'chip', '...');
      mirrorCard.appendChild(chips);
      appendText(mirrorCard, 'div', 'hint', 'Mirror targets are shown as shadow traffic branches in topology.');
      bindPreviewFocus(mirrorCard, policy, rule);
      root.appendChild(mirrorCard);
    }
  } else {
    const info = document.createElement('div');
    info.className = 'item';
    appendText(info, 'strong', '', 'Rule');
    appendText(info, 'div', 'hint', 'No rule selected in this policy.');
    root.appendChild(info);
  }
}

function renderContextBreadcrumb(){
  const root = el('contextBreadcrumb');
  if(!root) return;
  const policy = selectedPolicy();
  const rule = (policy?.rules || []).find(item => item.id === state.selectedRuleId);
  const parts = ['window ' + (state.statsWindow || '1h')];
  if(policy){ parts.push('policy ' + (policy.display_name || policy.key_id)); } else { parts.push('no policy selected'); }
  if(rule){ parts.push('rule ' + (rule.id || '-')); }
  if(state.focusedPreviewLabel){ parts.push((state.focusedPreviewType || 'preview') + ' ' + state.focusedPreviewLabel); }
  root.textContent = parts.join(' | ');
}

function renderStatsWindowButtons(){
  const current = state.statsWindow || '1h';
  [['statsWindow5mBtn','5m'],['statsWindow1hBtn','1h'],['statsWindow24hBtn','24h']].forEach(([id, value]) => {
    const button = el(id);
    if(!button) return;
    button.className = current === value ? 'warn' : 'alt';
    button.title = current === value ? 'Current stats window.' : 'Switch stats window to ' + value + '.';
  });
}

function renderUsage(){
  const root = el('usageStats'); clearNode(root);
  const rows = state.usage || [];
  const windowHits = currentWindowHitMaps();
  renderStatsWindowButtons();
  renderContextBreadcrumb();
  renderCurrentNodeDetail();
  if(!rows.length){ appendStat(root, 'No usage data yet.', '0'); renderGlobalTopology(); return; }
  rows.forEach(row => {
    const node = document.createElement('div'); node.className='stat';
    appendText(node, 'span', '', row.display_name || row.key_id);
    appendText(node, 'b', '', row.requests_today);
    appendText(node, 'div', 'hint', 'minute ' + row.requests_minute + ' | inflight ' + row.inflight + ' | window ' + windowHits.label);
    root.appendChild(node);
  });
  renderGlobalTopology();
}
function renderHealth(){
  const root = el('healthStats');
  const warningsRoot = el('healthWarnings');
  clearNode(root);
  clearNode(warningsRoot);
  if(!root || !warningsRoot) return;
  const health = state.health || {};
  const counters = health.counters || {};
  const persistence = health.persistence || {};
  const security = health.security || {};
  const counts = health.counts || {};
  appendStat(root, 'Status', health.status || 'unknown');
  appendStat(root, 'Counters', counters.backend || 'local');
  appendStat(root, 'Redis', counters.redis_status || 'not_configured');
  appendStat(root, 'Failure Mode', counters.failure_mode || '-');
  appendStat(root, 'Redis Required', String(Boolean(counters.redis_required)));
  appendStat(root, 'Persistence', persistence.storage || 'memory');
  appendStat(root, 'Security', security.management_auth_enabled ? 'management tokens' : (security.ui_auth_enabled ? 'ui token only' : 'host-only'));
  appendStat(root, 'Policies', counts.key_policies || 0);
  const warnings = Array.isArray(health.warnings) ? health.warnings : [];
  if(!warnings.length){
    warningsRoot.appendChild(emptyNode('item', 'No diagnostics warnings.'));
    return;
  }
  warnings.forEach(warning => warningsRoot.appendChild(emptyNode('item', warning)));
}
function renderAudit(){
  const root = el('auditList'); clearNode(root);
  const rows = state.audit || [];
  if(!rows.length){ root.appendChild(emptyNode('item', 'No audit records yet.')); return; }
  rows.forEach(row => {
    const node = document.createElement('div'); node.className='item';
    appendText(node, 'strong', '', row.decision || 'pass');
    appendText(node, 'div', 'hint', (row.rule_id || '-') + ' | ' + (row.reason || '-'));
    appendText(node, 'div', 'hint', (row.requested_model || '-') + ' -> ' + (row.final_model || '-'));
    appendText(node, 'div', 'hint', 'event ' + (row.event_type || 'request') + ' | operator ' + (row.operator_action || '-') + ' | member ' + (row.target_member || '-'));
    appendText(node, 'div', 'hint', 'before ' + (row.before_state || '-'));
    appendText(node, 'div', 'hint', 'after ' + (row.after_state || '-'));
    node.addEventListener('click', async () => {
      try {
        const detail = await readJSON(api.auditDetail + '?time=' + encodeURIComponent(row.time || '') + '&rule=' + encodeURIComponent(row.rule_id || '') + '&reason=' + encodeURIComponent(row.reason || ''));
        el('auditDetailOut').textContent = JSON.stringify(detail, null, 2);
        const stats = el('auditDetailStats');
        clearNode(stats);
        ['decision','rule_id','reason','requested_model','final_model','event_type','operator_action','target_member','secondary','before_state','after_state'].forEach(key => {
          const value = detail[key];
          if(!value) return;
          appendStat(stats, key, value);
        });
        if(Array.isArray(detail.diff) && detail.diff.length){
          appendStat(stats, 'diff', detail.diff.length + ' changes');
        }
        log('Audit detail loaded.');
      } catch(err){ log(err.message, 'bad'); }
    });
    root.appendChild(node);
  });
}
function renderAuditSummary(){
  const root = el('auditSummaryStats'); clearNode(root);
  const groups = state.auditSummary?.total_by_decision || {};
  const byPolicy = state.auditSummary?.total_by_policy || {};
  const byModel = state.auditSummary?.total_by_model || {};
  if(!Object.keys(groups).length && !Object.keys(byPolicy).length && !Object.keys(byModel).length){ appendStat(root, 'No audit summary yet.', '0'); return; }
  Object.entries(groups).forEach(([key, value]) => {
    const node = appendStat(root, 'decision: ' + key, value);
    node.addEventListener('click', () => { el('auditDecisionFilter').value = key; refreshAll(); });
  });
  Object.entries(byPolicy).slice(0, 3).forEach(([key, value]) => {
    appendStat(root, 'policy: ' + key, value);
  });
  Object.entries(byModel).slice(0, 3).forEach(([key, value]) => {
    appendStat(root, 'final model: ' + key, value);
  });
  const timelineRoot = el('auditTimelineStats');
  if(timelineRoot){
    clearNode(timelineRoot);
    (state.auditSummary?.timeline || []).slice(-8).forEach(item => {
      appendStat(timelineRoot, item.window, item.count);
    });
  }
}
function renderTemplates(){
  const root = el('templateList'); clearNode(root);
  const search = el('templateSearchInput')?.value?.trim().toLowerCase() || '';
  const category = el('templateCategoryFilter')?.value || '';
  const scenario = el('templateScenarioFilter')?.value || '';
  const maturity = el('templateMaturityFilter')?.value || '';
  const tag = (el('templateTagFilter')?.value || '').trim().toLowerCase();
  const items = (state.templates || []).filter(item => {
    const matchesSearch = !search || (item.name || '').toLowerCase().includes(search) || (item.description || '').toLowerCase().includes(search) || (item.scenario || '').toLowerCase().includes(search) || (item.tags || []).join(' ').toLowerCase().includes(search);
    const matchesCategory = !category || item.category === category;
    const matchesScenario = !scenario || item.scenario === scenario;
    const matchesMaturity = !maturity || item.maturity === maturity;
    const matchesTag = !tag || (item.tags || []).some(value => (value || '').toLowerCase().includes(tag));
    return matchesSearch && matchesCategory && matchesScenario && matchesMaturity && matchesTag;
  });
  if(!items.length){ root.appendChild(emptyNode('item', 'No templates loaded.')); return; }
  items.forEach(item => {
    const node = document.createElement('div'); node.className='item';
    appendText(node, 'strong', '', item.name || item.id || 'Template');
    appendText(node, 'div', 'hint', item.description || '');
    appendText(node, 'div', 'hint', 'category: ' + (item.category || 'custom') + ' | scenario: ' + (item.scenario || '-') + ' | maturity: ' + (item.maturity || '-'));
    const chips = document.createElement('div');
    chips.className = 'chips';
    const tags = item.tags || [];
    if(tags.length){
      tags.forEach(tagValue => appendText(chips, 'span', 'chip', tagValue));
    } else {
      appendText(chips, 'span', 'chip', 'no-tags');
    }
    node.appendChild(chips);
    appendText(node, 'div', 'hint', 'id: ' + (item.id || '-'));
    node.addEventListener('click', () => { const cloned = JSON.parse(JSON.stringify(item.rule || {})); if(!cloned.id){ cloned.id = 'template-rule-' + Date.now(); } appendRuleTemplate(cloned); hydrateRule(cloned); log('Template appended: ' + (item.name || item.id || 'template')); });
    const editBtn = document.createElement('button'); editBtn.className='alt'; editBtn.textContent='Edit'; editBtn.addEventListener('click', async (event) => { event.stopPropagation(); const name = prompt('Template name', item.name); if(name === null) return; const description = prompt('Template description', item.description || ''); if(description === null) return; const scenario = prompt('Template scenario', item.scenario || ''); if(scenario === null) return; const maturity = prompt('Template maturity', item.maturity || 'stable'); if(maturity === null) return; const tags = prompt('Template tags (comma separated)', (item.tags || []).join(',')); if(tags === null) return; try { await readJSON(api.templates + '?template_id=' + encodeURIComponent(item.id), {method:'PATCH', headers:{'Content-Type':'application/json'}, body: JSON.stringify({ ...item, name, description, scenario, maturity, tags: parseCSV(tags) })}); log('Template updated.'); await refreshAll(); } catch(err){ log(err.message, 'bad'); } });
    const cloneBtn = document.createElement('button'); cloneBtn.className='warn'; cloneBtn.textContent='Clone'; cloneBtn.addEventListener('click', async (event) => { event.stopPropagation(); try { await readJSON(api.templates + '/clone?template_id=' + encodeURIComponent(item.id), {method:'POST'}); log('Template cloned.'); await refreshAll(); } catch(err){ log(err.message, 'bad'); } });
    const useBtn = document.createElement('button'); useBtn.className='alt'; useBtn.textContent='Use'; useBtn.addEventListener('click', (event) => { event.stopPropagation(); appendRuleTemplate(item.rule || {}); log('Template inserted into current policy: ' + (item.name || item.id || 'template')); });
    const deleteBtn = document.createElement('button'); deleteBtn.className='danger'; deleteBtn.textContent='Delete'; deleteBtn.addEventListener('click', async (event) => { event.stopPropagation(); try { await readJSON(api.templates + '?template_id=' + encodeURIComponent(item.id), {method:'DELETE'}); log('Template deleted.'); await refreshAll(); } catch(err){ log(err.message, 'bad'); } });
    node.appendChild(editBtn);
    node.appendChild(cloneBtn);
    node.appendChild(useBtn);
    node.appendChild(deleteBtn);
    root.appendChild(node);
  });
}
function renderFallbackPreview(){
  const raw = el('ruleFallbackModelInput').value.trim();
  const chain = raw ? raw.split(',').map(v => v.trim()).filter(Boolean) : [];
  el('fallbackPreview').textContent = JSON.stringify(chain, null, 2);
  renderPoolHealth();
  renderRouteGraph();
}

function memberLabel(member){ return member?.model || ((member?.provider || '-') + '/' + (member?.suffix || '-')); }
function weightedRoutes(){ return Array.isArray(state.weightedRoutes) ? state.weightedRoutes : []; }
function setWeightedRoutes(items){ state.weightedRoutes = Array.isArray(items) ? items.slice() : []; renderWeightedRoutes(); }
function renderWeightedRoutes(){
  const root = el('weightedRoutesList');
  if(!root) return;
  clearNode(root);
  const items = weightedRoutes();
  if(!items.length){
    root.appendChild(textNode('div', 'item', 'No weighted routes configured.'));
    renderPoolHealth();
    renderRouteGraph();
    return;
  }
  items.forEach((member, index) => {
    const node = document.createElement('div');
    node.className = 'tableitem';
    const summary = document.createElement('div');
    appendText(summary, 'b', '', memberLabel(member));
    appendText(summary, 'span', '', 'status ' + (member.status || 'active') + ' | enabled ' + String(member.enabled !== false) + ' | priority ' + (member.priority ?? 100));
    node.appendChild(summary);
    appendText(node, 'span', '', 'weight ' + (member.weight ?? 1) + ' | health ' + (member.health ?? 100) + ' | cap ' + (member.traffic_cap ?? 100));
    const actions = document.createElement('div');
    actions.className = 'actions';
    appendButton(actions, 'alt', 'Edit', () => {
      el('weightedRouteModelInput').value = member.model || '';
      el('weightedRouteProviderInput').value = member.provider || '';
      el('weightedRouteSuffixInput').value = member.suffix || '';
      el('weightedRouteWeightInput').value = member.weight ?? 1;
      el('weightedRoutePriorityInput').value = member.priority ?? 100;
      el('weightedRouteEnabledInput').value = String(member.enabled !== false);
      el('weightedRouteStatusInput').value = member.status || 'active';
      el('weightedRouteReasonInput').value = member.reason || '';
      el('weightedRouteHealthInput').value = member.health ?? 100;
      el('weightedRouteTrafficCapInput').value = member.traffic_cap ?? 100;
      state.weightedRoutes.splice(index, 1);
      renderWeightedRoutes();
      log('Weighted route loaded into editor.');
    });
    appendButton(actions, 'danger', 'Delete', () => {
      state.weightedRoutes.splice(index, 1);
      renderWeightedRoutes();
      log('Weighted route removed.');
    });
    node.appendChild(actions);
    root.appendChild(node);
  });
  renderPoolHealth();
  renderRouteGraph();
}

function failoverHops(){ return Array.isArray(state.failoverHops) ? state.failoverHops : []; }
function setFailoverHops(items){ state.failoverHops = Array.isArray(items) ? items.slice() : []; renderFailoverHops(); }
function renderFailoverHops(){
  const root = el('failoverHopsList');
  if(!root) return;
  clearNode(root);
  const items = failoverHops();
  if(!items.length){ root.appendChild(textNode('div', 'item', 'No failover hops configured.')); return; }
  items.forEach((hop, index) => {
    const node = document.createElement('div'); node.className='tableitem';
    const label = hop.model || ((hop.provider || '-') + '/' + (hop.suffix || '-'));
    const summary = document.createElement('div');
    appendText(summary, 'b', '', label);
    appendText(summary, 'span', '', (hop.on_decision || 'reject') + ' | ' + (hop.reason || '-') + ' | enabled ' + String(hop.enabled !== false));
    node.appendChild(summary);
    appendText(node, 'span', '', 'hop ' + (index + 1));
    const actions = document.createElement('div'); actions.className='actions';
    const editBtn = document.createElement('button'); editBtn.className='alt'; editBtn.textContent='Edit'; editBtn.addEventListener('click', () => { el('failoverHopModelInput').value = hop.model || ''; el('failoverHopProviderInput').value = hop.provider || ''; el('failoverHopSuffixInput').value = hop.suffix || ''; el('failoverHopReasonInput').value = hop.reason || ''; el('failoverHopDecisionInput').value = hop.on_decision || 'reject'; el('failoverHopEnabledInput').value = String(hop.enabled !== false); state.failoverHops.splice(index, 1); renderFailoverHops(); log('Failover hop loaded into editor.'); });
    const deleteBtn = document.createElement('button'); deleteBtn.className='danger'; deleteBtn.textContent='Delete'; deleteBtn.addEventListener('click', () => { state.failoverHops.splice(index, 1); renderFailoverHops(); log('Failover hop removed.'); });
    actions.append(editBtn, deleteBtn); node.appendChild(actions); root.appendChild(node);
  });
  el('ruleFailoverHopsInput').value = JSON.stringify(items, null, 2);
}

function parseConditionGroups(inputID){
  const raw = el(inputID).value.trim();
  if(!raw) return [];
  const parsed = JSON.parse(raw);
  if(!Array.isArray(parsed)){ throw new Error(inputID + ' must be a JSON array.'); }
  return parsed;
}
function renderConditionGroups(){
  const root = el('conditionGroupsPreview');
  if(!root) return;
  clearNode(root);
  let anyOf = [];
  let allOf = [];
  try {
    anyOf = el('ruleUseAnyOf').value === 'true' ? parseConditionGroups('ruleAnyOfInput') : [];
    allOf = el('ruleUseAllOf').value === 'true' ? parseConditionGroups('ruleAllOfInput') : [];
  } catch(err) {
    const node = textNode('div', 'item', 'Invalid condition group JSON: ' + err.message);
    root.appendChild(node);
    return;
  }
  const addGroup = (kind, group, index) => {
    const node = document.createElement('div');
    node.className = 'item';
    appendText(node, 'strong', '', kind + ' #' + (index + 1));
    appendText(node, 'div', 'hint', 'models ' + ((group.models || []).length) + ' | paths ' + ((group.paths || []).length) + ' | providers ' + ((group.providers || []).length));
    appendText(node, 'div', 'hint', 'headers ' + Object.keys(group.headers || {}).length + ' | query ' + Object.keys(group.query || {}).length + ' | body ' + Object.keys(group.body_contains || {}).length + ' | metadata ' + Object.keys(group.metadata_contains || {}).length);
    const actions = document.createElement('div');
    actions.className = 'actions';
    appendButton(actions, 'alt', 'Load', (event) => {
      event.stopPropagation();
      loadMatchGroupIntoForm(group);
      log('Loaded ' + kind + ' condition group #' + (index + 1) + '.');
    });
    node.appendChild(actions);
    root.appendChild(node);
  };
  anyOf.forEach((group, index) => addGroup('Any-Of', group || {}, index));
  allOf.forEach((group, index) => addGroup('All-Of', group || {}, index));
  if(!anyOf.length && !allOf.length){
    root.appendChild(textNode('div', 'item', 'No condition groups configured.'));
  }
}

function bindTopologyActions(scope, policy, stage, rule, stageKey){
  if(!scope) return;
  scope.querySelectorAll('[data-topology-action]').forEach(button => {
    const action = button.getAttribute('data-topology-action');
    if(action === 'focus-policy'){
      button.onclick = (event) => { event.stopPropagation(); state.selectedKeyId = policy?.key_id || ''; state.selectedRuleId = ''; syncURLState(); hydrateSelectedPolicy(); renderKeys(); renderRules(); renderRouteGraph(); renderCurrentNodeDetail(); log('Focused policy from topology detail.'); };
    }
    if(action === 'copy-policy-summary'){
      button.onclick = async (event) => {
        event.stopPropagation();
        try {
          const summary = { policy: policy?.display_name || policy?.key_id || '', rules: (policy?.rules || []).length, limits: policy?.limits || {}, stage_policy: policy?.stage_policy || {} };
          await navigator.clipboard.writeText(JSON.stringify(summary, null, 2));
          log('Copied policy summary from topology detail.');
        } catch(err){ log(err.message || 'Clipboard copy failed.', 'bad'); }
      };
    }
    if(action === 'toggle-stage'){
      button.onclick = (event) => { event.stopPropagation(); if(stageKey) toggleTopology(stageKey); };
    }
    if(action === 'copy-stage-summary'){
      button.onclick = async (event) => {
        event.stopPropagation();
        try {
          const summary = { stage: stage || '', stage_hits_total: state.stageHitCounts[stage || ''] || 0, stage_hits_window: currentWindowHitMaps().stage[stage || ''] || 0, rules: (policy?.rules || []).filter(item => (item.stage || 'pre-check') === stage).map(item => item.id) };
          await navigator.clipboard.writeText(JSON.stringify(summary, null, 2));
          log('Copied stage summary from topology detail.');
        } catch(err){ log(err.message || 'Clipboard copy failed.', 'bad'); }
      };
    }
    if(action === 'focus-rule'){
      button.onclick = (event) => { event.stopPropagation(); if(!rule) return; state.selectedKeyId = policy?.key_id || ''; state.selectedRuleId = rule.id || ''; syncURLState(); hydrateSelectedPolicy(); hydrateRule(rule); renderKeys(); renderRules(); renderRouteGraph(); renderCurrentNodeDetail(); log('Focused rule from topology detail.'); };
    }
    if(action === 'open-rule-dry-run'){
      button.onclick = (event) => { event.stopPropagation(); if(!rule) return; applyDryRunPayloadFromRule(policy, rule); log('Prepared dry-run from topology detail.'); };
    }
    if(action === 'copy-rule-summary-inline'){
      button.onclick = async (event) => {
        event.stopPropagation();
        try {
          const summary = { rule: rule?.id || '', stage: rule?.stage || '', model: (rule?.match?.models || [])[0] || '', route_pool: rule?.actions?.route_pool?.name || '', failover_hops: (rule?.actions?.failover_hops || []).length, mirrors: (rule?.actions?.mirror_models || []).length };
          await navigator.clipboard.writeText(JSON.stringify(summary, null, 2));
          log('Copied rule summary from topology detail.');
        } catch(err){ log(err.message || 'Clipboard copy failed.', 'bad'); }
      };
    }
  });
}

function renderGlobalTopology(){
  const root = el('globalTopologyOut');
  if(!root) return;
  clearNode(root);
  root.className = 'topology-flow';
  const policies = state.policies?.key_policies || [];
  if(!policies.length){ root.appendChild(emptyNode('item', 'No policy topology yet.')); return; }
  policies.slice(0, 8).forEach(policy => {
    const rules = policy.rules || [];
    const pools = rules.filter(rule => rule.actions?.route_pool).length;
    const failovers = rules.filter(rule => (rule.actions?.failover_hops || []).length || (rule.actions?.failover_chain || []).length).length;
    const mirrors = rules.filter(rule => (rule.actions?.mirror_models || []).length).length;
    const totalRuleHits = rules.reduce((sum, rule) => sum + (state.ruleHitCounts[rule.id || ''] || 0), 0);
    const totalRuleHits5m = rules.reduce((sum, rule) => sum + (state.ruleHitCountsLast5m[rule.id || ''] || 0), 0);
    const totalRuleHitsHour = rules.reduce((sum, rule) => sum + (state.ruleHitCountsLastHour[rule.id || ''] || 0), 0);
    const totalRuleHits24h = rules.reduce((sum, rule) => sum + (state.ruleHitCountsLast24h[rule.id || ''] || 0), 0);
    const routeStageHits = state.stageHitCounts['route'] || 0;
    const windowHits = currentWindowHitMaps();
    const policyKey = 'policy:' + (policy.key_id || policy.display_name || 'unknown');
    const collapsed = topologyCollapsed(policyKey);
    const policyWrap = document.createElement('div');
    policyWrap.className = 'topology-policy' + (state.selectedKeyId === (policy.key_id || '') ? ' active' : '');
    policyWrap.title = 'Policy node: click to focus, double-click to collapse or expand. Rules: ' + rules.length + ', pools: ' + pools + ', failovers: ' + failovers + ', mirrors: ' + mirrors + '.';
    const policyNode = document.createElement('div');
    policyNode.className = 'item';
    appendText(policyNode, 'strong', '', policy.display_name || policy.key_id);
    appendText(policyNode, 'div', 'hint', 'rules ' + rules.length + ' | pools ' + pools + ' | failovers ' + failovers + ' | mirrors ' + mirrors);
    appendText(policyNode, 'div', statTone(totalRuleHitsHour), 'hits total ' + totalRuleHits + ' | 5m ' + totalRuleHits5m + ' | 1h ' + totalRuleHitsHour + ' | 24h ' + totalRuleHits24h + ' | current ' + windowHits.label);
    appendText(policyNode, 'div', 'hint', 'stage modes: pre-check ' + (policy.stage_policy?.['pre-check']?.mode || 'first-match') + ' -> route ' + (policy.stage_policy?.route?.mode || 'first-match') + ' | route-stage hits ' + routeStageHits);
    appendText(policyNode, 'div', 'hint', collapsed ? '[+] expand policy' : '[-] collapse policy');
    appendNodeDetail(policyNode, 'Policy Detail', [
      'Key id: ' + (policy.key_id || '-'),
      'Display: ' + (policy.display_name || '-'),
      'Limits: day ' + (policy.limits?.requests_per_day || 0) + ', min ' + (policy.limits?.requests_per_minute || 0) + ', inflight ' + (policy.limits?.max_inflight || 0)
    ]);
    policyNode.title = 'Policy summary: ' + (policy.display_name || policy.key_id) + ' | rules ' + rules.length + ' | pools ' + pools + ' | failovers ' + failovers + ' | mirrors ' + mirrors;
    policyNode.addEventListener('click', () => { state.selectedKeyId = policy.key_id || ''; state.selectedRuleId = ''; hydrateSelectedPolicy(); renderKeys(); renderRules(); renderRouteGraph(); renderCurrentNodeDetail(); log('Topology jumped to policy: ' + (policy.display_name || policy.key_id)); });
    bindTopologyActions(policyNode, policy, '', null, '');
    policyNode.addEventListener('dblclick', (event) => { event.preventDefault(); toggleTopology(policyKey); });
    policyWrap.appendChild(policyNode);
    if(!collapsed){
      const grouped = {};
      rules.forEach(rule => { const stage = rule.stage || 'pre-check'; if(!grouped[stage]) grouped[stage] = []; grouped[stage].push(rule); });
      ['pre-check','rewrite','route','mirror','post-audit'].forEach(stage => {
        const stageRules = grouped[stage] || [];
        if(!stageRules.length) return;
        const stageKey = policyKey + ':stage:' + stage;
        const stageCollapsed = topologyCollapsed(stageKey);
        const stageWrap = document.createElement('div');
        stageWrap.className = 'topology-stage' + (state.selectedRuleId && stageRules.some(rule => rule.id === state.selectedRuleId) ? ' active' : '');
        stageWrap.title = 'Stage node: double-click to collapse or expand this stage. Rules: ' + stageRules.length + ', total hits: ' + (state.stageHitCounts[stage] || 0) + ', current window hits: ' + (windowHits.stage[stage] || 0) + '.';
        const stageNode = document.createElement('div');
        stageNode.className = 'item';
        appendText(stageNode, 'strong', '', 'Stage: ' + stage);
        appendText(stageNode, 'div', statTone(state.stageHitCountsLastHour[stage] || 0), 'rules ' + stageRules.length + ' | hits ' + (state.stageHitCounts[stage] || 0) + ' | 5m ' + (state.stageHitCountsLast5m[stage] || 0) + ' | 1h ' + (state.stageHitCountsLastHour[stage] || 0) + ' | 24h ' + (state.stageHitCountsLast24h[stage] || 0));
        appendText(stageNode, 'div', 'hint', stageCollapsed ? '[+] expand stage' : '[-] collapse stage');
        const stageDetail = appendNodeDetail(stageNode, 'Stage Detail', [
          'Current window hits: ' + (windowHits.stage[stage] || 0),
          'Stage mode: ' + ((policy.stage_policy?.[stage]?.mode) || (stage === 'mirror' || stage === 'post-audit' ? 'continue-all' : 'first-match')),
          'Actions'
        ]);
        const stageActions = document.createElement('div');
        stageActions.className = 'actions';
        appendButton(stageActions, 'alt', 'Toggle Stage').dataset.topologyAction = 'toggle-stage';
        appendButton(stageActions, 'alt', 'Copy Summary').dataset.topologyAction = 'copy-stage-summary';
        stageDetail.appendChild(stageActions);
        stageNode.title = 'Stage ' + stage + ' | total hits ' + (state.stageHitCounts[stage] || 0) + ' | current window ' + (windowHits.stage[stage] || 0);
        stageNode.addEventListener('dblclick', (event) => { event.preventDefault(); toggleTopology(stageKey); });
        bindTopologyActions(stageNode, policy, stage, null, stageKey);
        stageWrap.appendChild(stageNode);
        if(!stageCollapsed){
          const arrow = document.createElement('div');
          arrow.className = 'flow-arrow';
          arrow.textContent = 'v';
          stageWrap.appendChild(arrow);
          stageRules.slice(0, 4).forEach(rule => {
            const child = document.createElement('div');
            const poolMembers = rule.actions?.route_pool?.members || rule.actions?.weighted_routes || [];
            const failoverItems = rule.actions?.failover_hops || [];
            const mirrorItems = (rule.actions?.mirror_models || []).map(model => ({ model }));
            const previewFocused = poolMembers.some(member => previewMatches('pool-member', previewChipLabel('pool-member', member))) || failoverItems.some(hop => previewMatches('failover-hop', previewChipLabel('failover-hop', hop))) || mirrorItems.some(item => previewMatches('mirror-target', previewChipLabel('mirror-target', item)));
            const previewLabel = previewFocused ? state.focusedPreviewLabel : '';
            child.className = 'topology-rule' + (state.selectedRuleId === (rule.id || '') ? ' active' : '') + (previewFocused ? ' preview-focus' : '');
            const ruleItem = document.createElement('div');
            ruleItem.className = 'item';
            appendText(ruleItem, 'strong', '', 'Rule: ' + (rule.id || '-'));
            appendText(ruleItem, 'div', statTone(state.ruleHitCountsLastHour[rule.id || ''] || 0), 'stage ' + stage + ' | total ' + (state.ruleHitCounts[rule.id || ''] || 0) + ' | 5m ' + (state.ruleHitCountsLast5m[rule.id || ''] || 0) + ' | 1h ' + (state.ruleHitCountsLastHour[rule.id || ''] || 0) + ' | 24h ' + (state.ruleHitCountsLast24h[rule.id || ''] || 0));
            const detailLines = [
              'Model: ' + ((rule.match?.models || [])[0] || '-'),
              'Route pool: ' + (rule.actions?.route_pool?.name || '-'),
              'Failover hops: ' + ((rule.actions?.failover_hops || []).length) + ' | Mirrors: ' + ((rule.actions?.mirror_models || []).length)
            ];
            if(previewLabel) detailLines.push('Preview focus: ' + previewLabel);
            detailLines.push('Targets');
            const ruleDetail = appendNodeDetail(ruleItem, 'Rule Detail', detailLines);
            appendSubnodeList(ruleDetail, 'pool-member', poolMembers, { kind: 'route', emptyText: 'No route pool members.', summary: (member, index) => 'member ' + (index + 1) + ' | weight ' + (member.weight ?? 1) + ' | priority ' + (member.priority ?? 100), meta: member => 'status ' + (member.status || 'active') + ' | health ' + (member.health ?? 100) + ' | cap ' + (member.traffic_cap ?? 100) });
            appendSubnodeList(ruleDetail, 'failover-hop', failoverItems, { kind: 'failover', emptyText: '', summary: (hop, index) => 'hop ' + (index + 1) + ' | decision ' + (hop.on_decision || 'reject'), meta: hop => (hop.reason || 'fallback') + ' | enabled ' + String(hop.enabled !== false) });
            appendSubnodeList(ruleDetail, 'mirror-target', mirrorItems, { kind: 'mirror', emptyText: '', summary: (_item, index) => 'mirror ' + (index + 1), meta: () => 'shadow traffic target' });
            appendText(ruleDetail, 'span', '', 'Actions');
            const ruleActions = document.createElement('div');
            ruleActions.className = 'actions';
            appendButton(ruleActions, 'alt', 'Focus Rule').dataset.topologyAction = 'focus-rule';
            appendButton(ruleActions, 'alt', 'Dry-Run').dataset.topologyAction = 'open-rule-dry-run';
            appendButton(ruleActions, 'alt', 'Copy Summary').dataset.topologyAction = 'copy-rule-summary-inline';
            ruleDetail.appendChild(ruleActions);
            child.appendChild(ruleItem);
            child.title = 'Rule ' + (rule.id || '-') + ' | current window hits ' + (windowHits.rule[rule.id || ''] || 0) + ' | match model ' + ((rule.match?.models || [])[0] || '-') + ' | route pool ' + (rule.actions?.route_pool?.name || '-') + ' | failover hops ' + ((rule.actions?.failover_hops || []).length) + ' | click to focus';
            child.addEventListener('click', () => { state.selectedKeyId = policy.key_id || ''; state.selectedRuleId = rule.id || ''; hydrateSelectedPolicy(); hydrateRule(rule); renderKeys(); renderRules(); renderRouteGraph(); renderCurrentNodeDetail(); log('Topology jumped to rule: ' + (rule.id || '-')); });
            bindTopologyActions(child, policy, stage, rule, stageKey);
            bindPreviewFocus(child, policy, rule);
            stageWrap.appendChild(child);
          });
        }
        policyWrap.appendChild(stageWrap);
      });
    }
    root.appendChild(policyWrap);
  });
}

function renderPoolHealth(){
  const root = el('poolHealthOut');
  if(!root) return;
  clearNode(root);
  const members = weightedRoutes();
  const usage = state.usage || [];
  const audit = state.audit || [];
  if(!members.length){ root.appendChild(textNode('div', 'item', 'No pool members configured.')); return; }
  const latestUsage = usage[0] || null;
  members.forEach((member, index) => {
    const label = memberLabel(member);
    const enabled = member.enabled !== false;
    const windowHits = currentWindowHitMaps();
    const hitCount = state.memberHitCounts[String(label).toLowerCase()] || audit.filter(item => (item.final_model || '').toLowerCase() === String(label).toLowerCase()).length;
    const hitWindow = windowHits.member[String(label).toLowerCase()] || 0;
    const node = document.createElement('div');
    node.className = 'item';
    node.title = 'Pool member health comes from configured values plus recent hit counts in the selected window.';
    appendText(node, 'strong', '', label);
    appendText(node, 'div', 'hint', 'status ' + (member.status || 'active') + ' | enabled ' + enabled + ' | priority ' + (member.priority ?? 100));
    appendText(node, 'div', 'hint', 'health ' + (member.health ?? 100) + ' | traffic cap ' + (member.traffic_cap ?? 100) + ' | weight ' + (member.weight ?? 1));
    appendText(node, 'div', 'hint', 'hits ' + hitCount + ' | current ' + windowHits.label + ' ' + hitWindow + (latestUsage ? ' | inflight ' + latestUsage.inflight + ' | minute ' + latestUsage.requests_minute : ''));
    appendText(node, 'div', 'hint', member.reason || 'no member note');
    const actionRows = [
      [['focus', 'Focus', 'alt'], ['active', 'Active', 'alt'], ['drain', 'Drain', 'warn'], ['offline', 'Offline', 'danger'], ['cap-down', 'Cap -10', 'alt'], ['cap-up', 'Cap +10', 'alt']],
      [['weight-down', 'Weight -1', 'alt'], ['weight-up', 'Weight +1', 'alt'], ['priority-down', 'Priority -10', 'alt'], ['priority-up', 'Priority +10', 'alt'], ['health-down', 'Health -10', 'alt'], ['health-up', 'Health +10', 'alt']]
    ];
    actionRows.forEach(row => {
      const actions = document.createElement('div');
      actions.className = 'actions';
      row.forEach(([op, labelText, className]) => {
        const button = appendButton(actions, className, labelText);
        button.dataset.memberOp = op;
        button.dataset.memberIndex = String(index);
      });
      node.appendChild(actions);
    });
    root.appendChild(node);
  });
  root.querySelectorAll('[data-member-op]').forEach(button => {
    button.addEventListener('click', async (event) => {
      event.stopPropagation();
      const index = Number(button.getAttribute('data-member-index'));
      const op = button.getAttribute('data-member-op') || '';
      const items = weightedRoutes();
      if(Number.isNaN(index) || index < 0 || index >= items.length) return;
      const target = items[index];
      const label = target.model || ((target.provider || '-') + '/' + (target.suffix || '-'));
      if(op === 'focus'){
        const previewLabel = previewChipLabel('pool-member', target);
        setPreviewFocus('pool-member', previewLabel);
        renderPoolHealth();
        renderGlobalTopology();
        renderCurrentNodeDetail();
        log('Focused pool member ' + previewLabel + '.');
        return;
      }
      try {
        if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a rule before applying member operations.'); }
        await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, member: label, operation: op, delta: 10, reason: op === 'drain' ? 'manual-drain' : (op === 'offline' ? 'manual-offline' : '') });
        await refreshAll();
        const refreshedPolicy = selectedPolicy();
        const refreshedRule = (refreshedPolicy?.rules || []).find(item => item.id === state.selectedRuleId);
        if(refreshedRule){
          hydrateRule(refreshedRule);
        }
        log('Persisted pool member ' + label + ' via ' + op + '.');
      } catch(err){
        log(err.message || 'Member operation failed.', 'bad');
      }
    });
  });
}

function renderRouteGraph(){
  const root = el('routeGraphOut');
  if(!root) return;
  clearNode(root);
  try {
    const rule = ruleFromVisualForm();
    const poolMembers = rule.actions?.route_pool?.members || rule.actions?.weighted_routes || [];
    const failoverItems = rule.actions?.failover_hops || failoverHops() || [];
    const blocks = [
      { title: 'Ingress', body: (rule.match?.paths || []).join(',') || 'any path', tone: 'ready' },
      { title: 'Stage', body: rule.stage || 'pre-check', tone: 'active' },
      { title: 'Route Pool', body: rule.actions?.route_pool ? ((rule.actions.route_pool.name || 'pool') + ' | ' + (rule.actions.route_pool.mode || 'weighted')) : 'none', tone: poolMembers.length ? 'active' : 'idle', meta: 'Members ' + poolMembers.length + ' | affinity ' + (rule.actions?.route_pool?.provider_affinity || '-') },
      { title: 'Failover', body: 'chain ' + String((rule.actions?.failover_chain || []).length) + ' | hops ' + String(failoverItems.length), tone: (failoverItems.length || (rule.actions?.failover_chain || []).length) ? 'warn' : 'idle', meta: 'Secondary targets activate when the primary route is not acceptable.' },
      { title: 'Mirror', body: (rule.actions?.mirror_models || []).join(',') || 'none', tone: (rule.actions?.mirror_models || []).length ? 'active' : 'idle', meta: 'Shadow targets receive mirrored routing metadata.' }
    ];
    blocks.forEach((block, index) => {
      const node = document.createElement('div');
      node.className = 'item';
      const badge = block.tone === 'warn' ? 'warn' : (block.tone === 'active' ? 'pill' : 'hint');
      appendText(node, 'strong', '', (index + 1) + '. ' + block.title);
      appendText(node, 'div', badge, block.body);
      appendText(node, 'div', 'hint', block.meta || 'Request enters gateway rule and flows to the next node.');
      if(index < blocks.length - 1){
        appendText(node, 'div', 'hint', '|');
        appendText(node, 'div', 'hint', 'v');
      }
      root.appendChild(node);
    });
  } catch {
    root.appendChild(emptyNode('item', 'No route graph yet.'));
  }
}
function renderRules(){
  const root = el('rulesList'); clearNode(root);
  const rules = selectedPolicy()?.rules || [];
  if(!rules.length){ root.appendChild(emptyNode('item', 'No rules yet.')); return; }
  rules.forEach(rule => {
    const node = document.createElement('div'); node.className='item' + (state.selectedRuleId === rule.id ? ' active' : '');
    const action = rule.actions?.route_to_model || (rule.actions?.route_pool?.name || '') || (rule.actions?.fallback_models || [])[0] || (rule.actions?.deny ? 'deny' : 'pass');
    const anyCount = (rule.match?.any_of || []).length;
    const allCount = (rule.match?.all_of || []).length;
    appendText(node, 'strong', '', rule.id || 'rule');
    appendText(node, 'div', 'hint', 'priority ' + (rule.priority || 0) + ' | on_match ' + (rule.on_match || 'stop'));
    appendText(node, 'div', 'hint', 'action: ' + action);
    appendText(node, 'div', 'hint', 'groups: any=' + anyCount + ' / all=' + allCount);
    node.addEventListener('click', () => { state.selectedRuleId = rule.id; clearPreviewFocus(); syncURLState(); hydrateRule(rule); renderRules(); renderCurrentNodeDetail(); renderGlobalTopology(); });
    root.appendChild(node);
  });
}
function hydrateRule(rule){
  if(!rule){ return; }
  el('ruleIdInput').value = rule.id || '';
  el('ruleStageInput').value = rule.stage || 'pre-check';
  el('ruleMatchModelInput').value = (rule.match?.models || []).join(',');
  el('rulePriorityInput').value = rule.priority || 10;
  el('rulePathInput').value = (rule.match?.paths || []).join(',');
  el('ruleProviderInput').value = (rule.match?.providers || []).join(',');
  const headerEntries = Object.entries(rule.match?.headers || {});
  el('ruleHeaderInput').value = headerEntries.map(([k,v]) => k + ':' + v).join(';');
  const queryEntries = Object.entries(rule.match?.query || {});
  el('ruleQueryInput').value = queryEntries.map(([k,v]) => k + ':' + v).join(';');
  const bodyEntries = Object.entries(rule.match?.body_contains || {});
  el('ruleBodyContainsInput').value = bodyEntries.map(([k,v]) => k + ':' + v).join(';');
  const metaEntries = Object.entries(rule.match?.metadata_contains || {});
  el('ruleMetadataInput').value = metaEntries.map(([k,v]) => k + ':' + v).join(';');
  el('ruleRouteToModelInput').value = rule.actions?.route_to_model || '';
  el('ruleForceProviderInput').value = rule.actions?.force_provider_prefix || '';
  el('ruleShardByInput').value = rule.actions?.shard_by || '';
  el('ruleMirrorModelsInput').value = (rule.actions?.mirror_models || []).join(',');
  el('ruleFallbackModelInput').value = (rule.actions?.fallback_models || []).join(',');
  el('ruleFailoverChainInput').value = (rule.actions?.failover_chain || []).join(',');
  el('routePoolNameInput').value = rule.actions?.route_pool?.name || '';
  el('routePoolModeInput').value = rule.actions?.route_pool?.mode || 'weighted';
  el('routePoolAffinityInput').value = rule.actions?.route_pool?.provider_affinity || '';
  setWeightedRoutes(rule.actions?.route_pool?.members || rule.actions?.weighted_routes || []);
  setFailoverHops(rule.actions?.failover_hops || []);
  el('ruleFailoverHopsInput').value = JSON.stringify(rule.actions?.failover_hops || [], null, 2);
  el('ruleFailoverChainInput').value = (rule.actions?.failover_chain || []).join(',');
  el('routePoolNameInput').value = rule.actions?.route_pool?.name || '';
  el('routePoolModeInput').value = rule.actions?.route_pool?.mode || 'weighted';
  el('routePoolAffinityInput').value = rule.actions?.route_pool?.provider_affinity || '';
  renderFallbackPreview();
  el('ruleDenyProviderInput').value = (rule.actions?.deny ? ((rule.match?.providers || [])[0] || '') : '');
  const reasonEntries = Object.entries(rule.actions?.tag_metadata || {});
  el('ruleReasonTagInput').value = reasonEntries.length ? (reasonEntries[0][0] + ':' + reasonEntries[0][1]) : '';
  const anyOf = rule.match?.any_of || [];
  const allOf = rule.match?.all_of || [];
  el('ruleUseAnyOf').value = String(anyOf.length > 0);
  el('ruleUseAllOf').value = String(allOf.length > 0);
  el('ruleAnyOfInput').value = anyOf.length ? JSON.stringify(anyOf, null, 2) : '[]';
  el('ruleAllOfInput').value = allOf.length ? JSON.stringify(allOf, null, 2) : '[]';
  if(anyOf.length > 0 || allOf.length > 0){
    log('Loaded rule with ' + anyOf.length + ' any-of groups and ' + allOf.length + ' all-of groups.');
  }
}
function hydrateSelectedPolicy(){
  const policy = selectedPolicy();
  if(!policy){
    el('policyKeyId').value = '';
    el('displayName').value = '';
    el('matchApiKey').value = '';
    el('matchApiKey').placeholder = 'paste top-level api key';
    el('policyEnabled').value = 'true';
    el('requestsPerDay').value = '';
    el('requestsPerMin').value = '';
    el('maxInflight').value = '';
    el('notBefore').value = '';
    el('notAfter').value = '';
    el('rulesBox').value = '[]';
    return;
  }
  el('policyKeyId').value = policy.key_id || '';
  el('displayName').value = policy.display_name || '';
  el('matchApiKey').value = policy.match_api_key || '';
  el('matchApiKey').placeholder = policy.masked_key ? ('leave blank to keep ' + policy.masked_key) : 'paste top-level api key';
  el('policyEnabled').value = String(policy.enabled !== false);
  el('requestsPerDay').value = policy.limits?.requests_per_day || '';
  el('requestsPerMin').value = policy.limits?.requests_per_minute || '';
  el('maxInflight').value = policy.limits?.max_inflight || '';
  el('notBefore').value = policy.limits?.not_before || '';
  el('notAfter').value = policy.limits?.not_after || '';
  el('stagePolicyPreCheck').value = policy.stage_policy?.['pre-check']?.mode || 'first-match';
  el('stagePolicyRewrite').value = policy.stage_policy?.rewrite?.mode || 'first-match';
  el('stagePolicyRoute').value = policy.stage_policy?.route?.mode || 'first-match';
  el('stagePolicyMirror').value = policy.stage_policy?.mirror?.mode || 'continue-all';
  el('stagePolicyPostAudit').value = policy.stage_policy?.['post-audit']?.mode || 'continue-all';
  el('rulesBox').value = JSON.stringify(policy.rules || [], null, 2);
}

function policyPayloadFromForm(){
  return {
    key_id: el('policyKeyId').value,
    display_name: el('displayName').value,
    match_api_key: el('matchApiKey').value,
    enabled: el('policyEnabled').value === 'true',
    limits: {
      requests_per_day: Number(el('requestsPerDay').value || 0),
      requests_per_minute: Number(el('requestsPerMin').value || 0),
      max_inflight: Number(el('maxInflight').value || 0),
      not_before: el('notBefore').value,
      not_after: el('notAfter').value
    },
    stage_policy: {
      'pre-check': { mode: el('stagePolicyPreCheck').value },
      rewrite: { mode: el('stagePolicyRewrite').value },
      route: { mode: el('stagePolicyRoute').value },
      mirror: { mode: el('stagePolicyMirror').value },
      'post-audit': { mode: el('stagePolicyPostAudit').value }
    },
    rules: JSON.parse(el('rulesBox').value || '[]')
  };
}
function auditQuery(){
  const params = new URLSearchParams({ limit: '30' });
  if(el('auditDecisionFilter').value) params.set('decision', el('auditDecisionFilter').value);
  if(el('auditRuleFilter').value.trim()) params.set('rule', el('auditRuleFilter').value.trim());
  if(el('auditReasonFilter').value.trim()) params.set('reason', el('auditReasonFilter').value.trim());
  if(el('auditKeyFilter').value.trim()) params.set('key', el('auditKeyFilter').value.trim());
  if(el('auditPolicyFilter').value.trim()) params.set('policy', el('auditPolicyFilter').value.trim());
  if(el('auditModelFilter').value.trim()) params.set('model', el('auditModelFilter').value.trim());
  if(el('auditProviderFilter').value.trim()) params.set('provider', el('auditProviderFilter').value.trim());
  if(el('auditEventTypeFilter').value.trim()) params.set('event_type', el('auditEventTypeFilter').value.trim());
  if(el('auditOperatorFilter').value.trim()) params.set('operator', el('auditOperatorFilter').value.trim());
  if(el('auditMemberFilter').value.trim()) params.set('member', el('auditMemberFilter').value.trim());
  if(el('auditFromFilter').value.trim()) params.set('from', el('auditFromFilter').value.trim());
  if(el('auditToFilter').value.trim()) params.set('to', el('auditToFilter').value.trim());
  return params.toString();
}
async function refreshAll(){
  try{
    const health = await readJSON(api.health);
    const [keys, policies, usage, audit, auditSummary, templates] = await Promise.all([readJSON(api.keys), readJSON(api.policies), readJSON(api.usage), readJSON(api.audit + '?' + auditQuery()), readJSON(api.auditSummary + '?' + auditQuery()), readJSON(api.templates)]);
    state.keys = keys.keys || [];
    state.policies = policies || { key_policies: [], default_policy: {} };
    state.usage = usage.usage || [];
    state.memberHitCounts = usage.member_hits || {};
    state.ruleHitCounts = usage.rule_hits || {};
    state.stageHitCounts = usage.stage_hits || {};
    state.memberHitCountsLast5m = usage.member_hits_last_5m || {};
    state.ruleHitCountsLast5m = usage.rule_hits_last_5m || {};
    state.stageHitCountsLast5m = usage.stage_hits_last_5m || {};
    state.memberHitCountsLastHour = usage.member_hits_last_hour || {};
    state.ruleHitCountsLastHour = usage.rule_hits_last_hour || {};
    state.stageHitCountsLastHour = usage.stage_hits_last_hour || {};
    state.memberHitCountsLast24h = usage.member_hits_last_24h || {};
    state.ruleHitCountsLast24h = usage.rule_hits_last_24h || {};
    state.stageHitCountsLast24h = usage.stage_hits_last_24h || {};
    state.audit = audit.items || [];
    state.auditSummary = auditSummary || { total_by_decision: {}, total_by_reason: {}, total_by_rule: {}, total_by_policy: {}, total_by_model: {} };
    state.templates = templates.items || [];
    state.health = health || {};
    if(!state.selectedKeyId && state.keys.length){ state.selectedKeyId = state.keys[0].key_id; }
    syncURLState();
    if(state.selectedKeyId && !(state.keys || []).some(item => item.key_id === state.selectedKeyId)){ state.selectedKeyId = state.keys[0]?.key_id || ''; }
    el('policiesBox').value = JSON.stringify(state.policies, null, 2);
    renderKeys();
    renderUsage();
    renderHealth();
    renderAuditSummary();
    renderAudit();
    hydrateSelectedPolicy();
    renderRules();
    renderTemplates();
    renderFallbackPreview();
    renderWeightedRoutes();
    renderConditionGroups();
    renderFailoverHops();
    renderPoolHealth();
    renderRouteGraph();
    renderDryRunDecision(null);
    renderDryRunStageTrace([]);
    setManagementStatus('Connected');
    log('Gateway data refreshed.');
  }catch(err){ log(err.message, 'bad'); }
}
async function connectGateway(){
  const key = currentManagementKey();
  if(!key){
    setManagementStatus('Enter management key first', 'bad');
    log('Enter the CPAMC management key before connecting.', 'bad');
    return;
  }
  storeManagementKey(key);
  managementAuthBlocked = false;
  setManagementStatus('Connecting...');
  await refreshAll();
}
function initManagementAccess(){
  const input = el('managementKeyInput');
  const stored = loadManagementKey();
  if(input && stored) input.value = stored;
  el('saveManagementKeyBtn').addEventListener('click', () => {
    const key = currentManagementKey();
    if(!key){
      setManagementStatus('Enter management key first', 'bad');
      log('Enter the CPAMC management key before saving.', 'bad');
      return;
    }
    storeManagementKey(key);
    managementAuthBlocked = false;
    setManagementStatus('Saved for this session');
    log('Management key saved for this browser session.');
  });
  el('clearManagementKeyBtn').addEventListener('click', () => {
    storeManagementKey('');
    managementAuthBlocked = false;
    if(input) input.value = '';
    setManagementStatus('Not connected');
    log('Management key cleared from this browser session.');
  });
  el('connectManagementBtn').addEventListener('click', connectGateway);
  if(stored){
    setManagementStatus('Stored key found');
    connectGateway();
  }else{
    setManagementStatus('Not connected');
    log('Enter the CPAMC management key, then click Connect / Refresh.');
  }
}
el('refreshBtn').addEventListener('click', connectGateway);
el('statsWindow5mBtn').addEventListener('click', () => { state.statsWindow = '5m'; syncURLState(); renderUsage(); renderPoolHealth(); renderGlobalTopology(); });
el('statsWindow1hBtn').addEventListener('click', () => { state.statsWindow = '1h'; syncURLState(); renderUsage(); renderPoolHealth(); renderGlobalTopology(); });
el('statsWindow24hBtn').addEventListener('click', () => { state.statsWindow = '24h'; syncURLState(); renderUsage(); renderPoolHealth(); renderGlobalTopology(); });
el('saveAllBtn').addEventListener('click', async () => {
  try{
    const parsed = JSON.parse(el('policiesBox').value || '{}');
    await readJSON(api.policies, {method:'PUT', headers:{'Content-Type':'application/json'}, body: JSON.stringify(parsed)});
    log('Whole policy set saved.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('ruleFallbackModelInput').addEventListener('input', renderFallbackPreview);
el('ruleAnyOfInput').addEventListener('input', renderConditionGroups);
el('ruleAllOfInput').addEventListener('input', renderConditionGroups);
el('applyTemplateFilterBtn').addEventListener('click', renderTemplates);
el('exportTemplatesBtn').addEventListener('click', async () => {
  try{
    const data = await readJSON(api.exportTemplates);
    el('dryRunOut').textContent = JSON.stringify(data, null, 2);
    renderDryRunStageTrace([]);
    log('Templates exported to output panel.');
  }catch(err){ log(err.message, 'bad'); }
});
el('importTemplatesBtn').addEventListener('click', async () => {
  try{
    const raw = el('templateImportBox').value.trim();
    if(!raw){ throw new Error('Template import JSON is empty.'); }
    await readJSON(api.importTemplates, {method:'POST', headers:{'Content-Type':'application/json'}, body: raw});
    log('Templates imported.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('buildFallbackChainBtn').addEventListener('click', () => {
  const fallback = el('ruleFallbackModelInput').value.trim();
  const failover = el('ruleFailoverChainInput').value.trim();
  const failoverHopsRaw = el('ruleFailoverHopsInput').value.trim();
  const routePoolName = el('routePoolNameInput').value.trim();
  const routePoolMode = el('routePoolModeInput').value;
  const routePoolAffinity = el('routePoolAffinityInput').value.trim();
  if(!fallback){ log('Enter fallback models first.', 'bad'); return; }
  const rule = appendRuleTemplate({ id: 'fallback-chain-' + Date.now(), enabled: true, priority: Number(el('rulePriorityInput').value || 10), on_match: 'stop', match: { models: el('ruleMatchModelInput').value.trim() ? el('ruleMatchModelInput').value.trim().split(',').map(v => v.trim()).filter(Boolean) : [] }, actions: { fallback_models: fallback.split(',').map(v => v.trim()).filter(Boolean) } });
  hydrateRule(rule);
  log('Fallback chain rule appended.');
});
function fallbackChain(){ const raw = el('ruleFallbackModelInput').value.trim(); return raw ? raw.split(',').map(v => v.trim()).filter(Boolean) : []; }
function setFallbackChain(chain){ el('ruleFallbackModelInput').value = chain.join(','); renderFallbackPreview(); }
el('addFallbackHopBtn').addEventListener('click', () => { const next = prompt('Fallback model'); if(!next) return; const chain = fallbackChain(); chain.push(next.trim()); setFallbackChain(chain); log('Fallback hop added.'); });
el('removeFallbackHopBtn').addEventListener('click', () => { const chain = fallbackChain(); if(!chain.length) return; chain.pop(); setFallbackChain(chain); log('Fallback hop removed.'); });
el('sortFallbackHopsBtn').addEventListener('click', () => { const chain = fallbackChain().sort(); setFallbackChain(chain); log('Fallback chain sorted.'); });
el('addFailoverHopBtn').addEventListener('click', () => {
  const model = el('failoverHopModelInput').value.trim();
  const provider = el('failoverHopProviderInput').value.trim();
  const suffix = el('failoverHopSuffixInput').value.trim();
  if(!model && !(provider && suffix)){ log('Enter a failover model or provider+suffix first.', 'bad'); return; }
  state.failoverHops.push({ model, provider, suffix, reason: el('failoverHopReasonInput').value.trim(), on_decision: el('failoverHopDecisionInput').value || 'reject', enabled: el('failoverHopEnabledInput').value !== 'false' });
  el('failoverHopModelInput').value = '';
  el('failoverHopProviderInput').value = '';
  el('failoverHopSuffixInput').value = '';
  el('failoverHopReasonInput').value = '';
  el('failoverHopDecisionInput').value = 'reject';
  el('failoverHopEnabledInput').value = 'true';
  renderFailoverHops();
  renderRouteGraph();
  log('Failover hop added.');
});
el('clearFailoverHopsBtn').addEventListener('click', () => { state.failoverHops = []; renderFailoverHops(); renderRouteGraph(); log('Failover hops cleared.'); });
el('addWeightedRouteBtn').addEventListener('click', () => { const model = el('weightedRouteModelInput').value.trim(); const provider = el('weightedRouteProviderInput').value.trim(); const suffix = el('weightedRouteSuffixInput').value.trim(); if(!model && !(provider && suffix)){ log('Enter a weighted route model or provider+suffix first.', 'bad'); return; } const weight = Number(el('weightedRouteWeightInput').value || 1) || 1; const priority = Number(el('weightedRoutePriorityInput').value || 100) || 100; const enabled = el('weightedRouteEnabledInput').value === 'true'; const status = el('weightedRouteStatusInput').value || 'active'; const reason = el('weightedRouteReasonInput').value.trim(); const health = Number(el('weightedRouteHealthInput').value || 100) || 100; const trafficCap = Number(el('weightedRouteTrafficCapInput').value || 100) || 100; state.weightedRoutes.push({ model, provider, suffix, weight, priority, enabled, status, reason, health, traffic_cap: trafficCap }); el('weightedRouteModelInput').value=''; el('weightedRouteProviderInput').value=''; el('weightedRouteSuffixInput').value=''; el('weightedRouteWeightInput').value='1'; el('weightedRoutePriorityInput').value='100'; el('weightedRouteEnabledInput').value='true'; el('weightedRouteStatusInput').value='active'; el('weightedRouteReasonInput').value=''; el('weightedRouteHealthInput').value='100'; el('weightedRouteTrafficCapInput').value='100'; renderWeightedRoutes(); log('Weighted route added.'); });
el('clearWeightedRoutesBtn').addEventListener('click', () => { state.weightedRoutes = []; renderWeightedRoutes(); log('Weighted routes cleared.'); });
el('sortWeightedRoutesBtn').addEventListener('click', () => { state.weightedRoutes = weightedRoutes().slice().sort((a, b) => memberLabel(a).localeCompare(memberLabel(b))); renderWeightedRoutes(); log('Weighted routes sorted.'); });
el('saveTemplateBtn').addEventListener('click', async () => {
  try{
    const rule = ruleFromVisualForm();
    const payload = { id: 'template-' + Date.now(), name: rule.id || 'Custom Template', category: 'custom', scenario: 'custom', maturity: 'experimental', tags: ['custom'], description: 'Saved from current rule editor.', rule };
    await readJSON(api.templates, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    log('Template saved.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('newPolicyBtn').addEventListener('click', async () => {
  try{
    const payload = { display_name: 'New Policy', match_api_key: '', enabled: true, limits: { requests_per_day: 0, requests_per_minute: 0, max_inflight: 0 }, rules: [] };
    await readJSON(api.addPolicy, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    log('New key policy added. Fill in the API key and save it.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('clonePolicyBtn').addEventListener('click', async () => {
  try{
    if(!state.selectedKeyId){ throw new Error('Select a key policy first.'); }
    await readJSON(api.clonePolicy + '?key_id=' + encodeURIComponent(state.selectedKeyId), {method:'POST'});
    log('Selected policy cloned. New copy has no bound API key yet.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('exportPolicyBundleBtn').addEventListener('click', async () => {
  try{
    const name = encodeURIComponent(el('policyBundleName').value.trim() || 'gateway-policy-bundle');
    const out = await readJSON(api.exportPolicies + '?name=' + name);
    el('policiesBox').value = JSON.stringify(out, null, 2);
    log('Policy bundle exported into Advanced JSON panel.');
  }catch(err){ log(err.message, 'bad'); }
});
el('importPolicyBundleBtn').addEventListener('click', async () => {
  try{
    const raw = el('policiesBox').value.trim();
    if(!raw){ throw new Error('Policies JSON is empty.'); }
    const mode = el('policyBundleMode').value || 'merge';
    await readJSON(api.importPolicies + '?mode=' + encodeURIComponent(mode), {method:'POST', headers:{'Content-Type':'application/json'}, body: raw});
    log('Policy bundle imported with mode: ' + mode + '.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('exportAuditBtn').addEventListener('click', async () => {
  try{
    const audit = await readJSON(api.audit + '?limit=200&' + auditQuery());
    el('dryRunOut').textContent = JSON.stringify(audit, null, 2);
    renderDryRunStageTrace([]);
    log('Audit exported to output panel.');
  }catch(err){ log(err.message, 'bad'); }
});
el('applyAuditFilterBtn').addEventListener('click', refreshAll);
el('lastHourAuditBtn').addEventListener('click', () => { const now = new Date(); const from = new Date(now.getTime() - 60 * 60 * 1000); el('auditFromFilter').value = from.toISOString(); el('auditToFilter').value = now.toISOString(); refreshAll(); });
el('todayAuditBtn').addEventListener('click', () => { const now = new Date(); const from = new Date(now); from.setHours(0,0,0,0); el('auditFromFilter').value = from.toISOString(); el('auditToFilter').value = now.toISOString(); refreshAll(); });
el('savePolicyBtn').addEventListener('click', async () => {
  try{
    if(!state.selectedKeyId){ throw new Error('Select a key policy first.'); }
    const payload = policyPayloadFromForm();
    await readJSON(api.policies + '?key_id=' + encodeURIComponent(state.selectedKeyId), {method:'PATCH', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    log('Selected policy saved.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('deletePolicyBtn').addEventListener('click', async () => {
  try{
    if(!state.selectedKeyId){ throw new Error('Select a key policy first.'); }
    await readJSON(api.policies + '?key_id=' + encodeURIComponent(state.selectedKeyId), {method:'DELETE'});
    state.selectedKeyId = '';
    log('Selected policy deleted.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('resetUsageBtn').addEventListener('click', async () => {
  try{
    if(!state.selectedKeyId){ throw new Error('Select a key policy first.'); }
    await readJSON(api.resetUsage + '?key_id=' + encodeURIComponent(state.selectedKeyId), {method:'POST'});
    log('Selected usage reset.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
function ruleFromVisualForm(){
  const model = el('ruleMatchModelInput').value.trim();
  const matchPath = el('rulePathInput').value.trim();
  const matchProvider = el('ruleProviderInput').value.trim();
  const matchHeader = el('ruleHeaderInput').value.trim();
  const matchQuery = el('ruleQueryInput').value.trim();
  const bodyContains = el('ruleBodyContainsInput').value.trim();
  const metadataContains = el('ruleMetadataInput').value.trim();
  const routeTo = el('ruleRouteToModelInput').value.trim();
  const forceProvider = el('ruleForceProviderInput').value.trim();
  const mirrorModels = el('ruleMirrorModelsInput').value.trim();
  const shardBy = el('ruleShardByInput').value;
  const fallback = el('ruleFallbackModelInput').value.trim();
  const failover = el('ruleFailoverChainInput').value.trim();
  const failoverHopsRaw = el('ruleFailoverHopsInput').value.trim();
  const routePoolName = el('routePoolNameInput').value.trim();
  const routePoolMode = el('routePoolModeInput').value;
  const routePoolAffinity = el('routePoolAffinityInput').value.trim();
  const denyProvider = el('ruleDenyProviderInput').value.trim();
  const reasonTag = el('ruleReasonTagInput').value.trim();
  const useAnyOf = el('ruleUseAnyOf').value === 'true';
  const useAllOf = el('ruleUseAllOf').value === 'true';
  const anyOfRaw = el('ruleAnyOfInput').value.trim();
  const allOfRaw = el('ruleAllOfInput').value.trim();
  const rule = { id: el('ruleIdInput').value.trim() || ('rule-' + Date.now()), enabled: true, priority: Number(el('rulePriorityInput').value || 10), stage: el('ruleStageInput').value || 'pre-check', on_match: 'stop', match: {}, actions: {} };
  if(model){ rule.match.models = parseCSV(model); }
  if(matchPath){ rule.match.paths = parseCSV(matchPath); }
  if(matchProvider){ rule.match.providers = parseCSV(matchProvider); }
  if(matchHeader){ rule.match.headers = parsePairs(matchHeader); }
  if(matchQuery){ rule.match.query = parsePairs(matchQuery); }
  if(bodyContains){ rule.match.body_contains = parsePairs(bodyContains); }
  if(metadataContains){ rule.match.metadata_contains = parsePairs(metadataContains); }
  if(routeTo){ rule.actions.route_to_model = routeTo; }
  if(forceProvider){ rule.actions.force_provider_prefix = forceProvider; }
  if(mirrorModels){ rule.actions.mirror_models = parseCSV(mirrorModels); }
  if(shardBy){ rule.actions.shard_by = shardBy; }
  const routes = weightedRoutes();
  if(routes.length){
    if(routePoolName || routePoolAffinity){
      rule.actions.route_pool = { name: routePoolName, mode: routePoolMode || 'weighted', provider_affinity: routePoolAffinity, members: routes };
    } else {
      rule.actions.weighted_routes = routes;
    }
  }
  if(fallback){ rule.actions.fallback_models = parseCSV(fallback); }
  if(failover){ rule.actions.failover_chain = parseCSV(failover); }
  if(failoverHops().length){ rule.actions.failover_hops = failoverHops(); } else if(failoverHopsRaw){ rule.actions.failover_hops = JSON.parse(failoverHopsRaw); }
  if(denyProvider){ rule.match.providers = [denyProvider]; rule.actions.deny = { status_code: 403, message: 'provider denied', code: 'gateway_provider_denied' }; }
  if(reasonTag){ rule.actions.tag_metadata = { reason: reasonTag }; }
  if(useAnyOf){ rule.match.any_of = anyOfRaw ? JSON.parse(anyOfRaw) : [{ paths: rule.match.paths || [], models: rule.match.models || [] }]; }
  if(useAllOf){ rule.match.all_of = allOfRaw ? JSON.parse(allOfRaw) : [{ providers: rule.match.providers || [], query: rule.match.query || {} }]; }
  return rule;
}
function parsePairs(raw){
  const out = {};
  raw.split(';').map(v => v.trim()).filter(Boolean).forEach(item => { const parts = item.split(':'); if(parts.length >= 2){ out[parts.shift().trim()] = parts.join(':').trim(); } });
  return out;
}
function parseCSV(raw){ return String(raw || '').split(',').map(v => v.trim()).filter(Boolean); }
function appendRuleTemplate(template){
  const current = JSON.parse(el('rulesBox').value || '[]');
  const next = JSON.parse(JSON.stringify(template));
  if(!next.id){ next.id = 'rule-' + Date.now(); }
  current.push(next);
  el('rulesBox').value = JSON.stringify(current, null, 2);
  state.selectedRuleId = next.id || '';
  renderRules();
  return next;
}
el('addRouteToModelRuleBtn').addEventListener('click', () => { const rule = appendRuleTemplate({ id: 'route-' + Date.now(), enabled: true, priority: 10, stage: 'route', on_match: 'stop', match: { models: ['gpt-5.5'] }, actions: { route_to_model: 'openai/gpt-5.4' } }); hydrateRule(rule); });
el('addFallbackRuleBtn').addEventListener('click', () => { const rule = appendRuleTemplate({ id: 'fallback-' + Date.now(), enabled: true, priority: 20, stage: 'route', on_match: 'stop', match: { models: ['gpt-5.5'] }, actions: { fallback_models: ['openai/gpt-5.4-mini'] } }); hydrateRule(rule); });
el('addDenyRuleBtn').addEventListener('click', () => { const rule = appendRuleTemplate({ id: 'deny-' + Date.now(), enabled: true, priority: 30, stage: 'pre-check', on_match: 'stop', match: { providers: ['claude'] }, actions: { deny: { status_code: 403, message: 'provider denied', code: 'gateway_provider_denied' } } }); hydrateRule(rule); });
function currentMatchGroupFromForm(){
  const group = {};
  const model = el('ruleMatchModelInput').value.trim();
  const path = el('rulePathInput').value.trim();
  const provider = el('ruleProviderInput').value.trim();
  const header = el('ruleHeaderInput').value.trim();
  const query = el('ruleQueryInput').value.trim();
  const body = el('ruleBodyContainsInput').value.trim();
  const metadata = el('ruleMetadataInput').value.trim();
  if(model){ group.models = model.split(',').map(v => v.trim()).filter(Boolean); }
  if(path){ group.paths = path.split(',').map(v => v.trim()).filter(Boolean); }
  if(provider){ group.providers = provider.split(',').map(v => v.trim()).filter(Boolean); }
  if(header){ group.headers = parsePairs(header); }
  if(query){ group.query = parsePairs(query); }
  if(body){ group.body_contains = parsePairs(body); }
  if(metadata){ group.metadata_contains = parsePairs(metadata); }
  return group;
}
el('buildAnyOfBtn').addEventListener('click', () => {
  const next = JSON.parse(el('ruleAnyOfInput').value || '[]');
  next.push(currentMatchGroupFromForm());
  el('ruleAnyOfInput').value = JSON.stringify(next, null, 2);
  el('ruleUseAnyOf').value = 'true';
  log('Any-Of group appended from current inputs.');
});
el('buildAllOfBtn').addEventListener('click', () => {
  const next = JSON.parse(el('ruleAllOfInput').value || '[]');
  next.push(currentMatchGroupFromForm());
  el('ruleAllOfInput').value = JSON.stringify(next, null, 2);
  el('ruleUseAllOf').value = 'true';
  log('All-Of group appended from current inputs.');
});
el('clearConditionGroupsBtn').addEventListener('click', () => {
  el('ruleAnyOfInput').value = '[]';
  el('ruleAllOfInput').value = '[]';
  el('ruleUseAnyOf').value = 'false';
  el('ruleUseAllOf').value = 'false';
  log('Condition groups cleared.');
});
el('syncCurrentRuleToGroupsBtn').addEventListener('click', () => {
  const current = currentMatchGroupFromForm();
  if(el('ruleUseAnyOf').value === 'true'){
    el('ruleAnyOfInput').value = JSON.stringify([current], null, 2);
  }
  if(el('ruleUseAllOf').value === 'true'){
    el('ruleAllOfInput').value = JSON.stringify([current], null, 2);
  }
  if(el('ruleUseAnyOf').value !== 'true' && el('ruleUseAllOf').value !== 'true'){
    el('ruleAnyOfInput').value = JSON.stringify([current], null, 2);
    el('ruleUseAnyOf').value = 'true';
  }
  log('Current rule inputs synced into active condition groups.');
});
el('popConditionGroupBtn').addEventListener('click', () => {
  let changed = false;
  const anyOf = JSON.parse(el('ruleAnyOfInput').value || '[]');
  const allOf = JSON.parse(el('ruleAllOfInput').value || '[]');
  if(el('ruleUseAnyOf').value === 'true' && anyOf.length > 0){ anyOf.pop(); el('ruleAnyOfInput').value = JSON.stringify(anyOf, null, 2); changed = true; }
  if(el('ruleUseAllOf').value === 'true' && allOf.length > 0){ allOf.pop(); el('ruleAllOfInput').value = JSON.stringify(allOf, null, 2); changed = true; }
  if(changed){ log('Removed last condition group item.'); }
});
function loadMatchGroupIntoForm(group){
  if(!group) return;
  el('ruleMatchModelInput').value = (group.models || []).join(',');
  el('rulePathInput').value = (group.paths || []).join(',');
  el('ruleProviderInput').value = (group.providers || []).join(',');
  el('ruleHeaderInput').value = Object.entries(group.headers || {}).map(([k,v]) => k + ':' + v).join(';');
  el('ruleQueryInput').value = Object.entries(group.query || {}).map(([k,v]) => k + ':' + v).join(';');
  el('ruleBodyContainsInput').value = Object.entries(group.body_contains || {}).map(([k,v]) => k + ':' + v).join(';');
  el('ruleMetadataInput').value = Object.entries(group.metadata_contains || {}).map(([k,v]) => k + ':' + v).join(';');
}
el('loadFirstAnyOfBtn').addEventListener('click', () => {
  const groups = JSON.parse(el('ruleAnyOfInput').value || '[]');
  if(!groups.length){ log('No Any-Of group to load.', 'bad'); return; }
  loadMatchGroupIntoForm(groups[0]);
  log('Loaded first Any-Of group into form.');
});
el('loadFirstAllOfBtn').addEventListener('click', () => {
  const groups = JSON.parse(el('ruleAllOfInput').value || '[]');
  if(!groups.length){ log('No All-Of group to load.', 'bad'); return; }
  loadMatchGroupIntoForm(groups[0]);
  log('Loaded first All-Of group into form.');
});
el('pushRuleFromFormBtn').addEventListener('click', () => { const rule = appendRuleTemplate(ruleFromVisualForm()); hydrateRule(rule); });
el('saveRuleBtn').addEventListener('click', async () => {
  try{
    if(!state.selectedKeyId){ throw new Error('Select a key policy first.'); }
    const payload = ruleFromVisualForm();
    if(state.selectedRuleId){
      await readJSON(api.rule + '?key_id=' + encodeURIComponent(state.selectedKeyId) + '&rule_id=' + encodeURIComponent(state.selectedRuleId), {method:'PATCH', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      log('Selected rule saved.');
    } else {
      await readJSON(api.addRule + '?key_id=' + encodeURIComponent(state.selectedKeyId), {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      log('Rule added to selected policy.');
    }
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
el('deleteRuleBtn').addEventListener('click', async () => {
  try{
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a rule first.'); }
    await readJSON(api.rule + '?key_id=' + encodeURIComponent(state.selectedKeyId) + '&rule_id=' + encodeURIComponent(state.selectedRuleId), {method:'DELETE'});
    state.selectedRuleId = '';
    log('Selected rule deleted.');
    await refreshAll();
  }catch(err){ log(err.message, 'bad'); }
});
function currentRules(){ return JSON.parse(el('rulesBox').value || '[]'); }
function replaceRules(next){ el('rulesBox').value = JSON.stringify(next, null, 2); renderRules(); }
function selectedRuleIndex(){ return currentRules().findIndex(rule => rule.id === state.selectedRuleId); }
el('cloneRuleBtn').addEventListener('click', () => {
  const rules = currentRules();
  const idx = selectedRuleIndex();
  if(idx < 0){ log('Select a rule to clone.', 'bad'); return; }
  const clone = JSON.parse(JSON.stringify(rules[idx]));
  clone.id = clone.id + '-copy';
  rules.splice(idx + 1, 0, clone);
  replaceRules(rules);
  state.selectedRuleId = clone.id;
  hydrateRule(clone);
  renderRules();
  log('Rule cloned.');
});
el('moveRuleUpBtn').addEventListener('click', () => {
  const rules = currentRules();
  const idx = selectedRuleIndex();
  if(idx <= 0){ return; }
  const item = rules[idx];
  rules.splice(idx, 1);
  rules.splice(idx - 1, 0, item);
  replaceRules(rules);
  renderRules();
});
el('moveRuleDownBtn').addEventListener('click', () => {
  const rules = currentRules();
  const idx = selectedRuleIndex();
  if(idx < 0 || idx >= rules.length - 1){ return; }
  const item = rules[idx];
  rules.splice(idx, 1);
  rules.splice(idx + 1, 0, item);
  replaceRules(rules);
  renderRules();
});
el('dryRunBtn').addEventListener('click', async () => {
  try{
    const reqBody = JSON.parse(el('dryBody').value || '{}');
    const metadata = {'access.api_key': el('dryKey').value, 'access.key_id': state.selectedKeyId || '', 'request_path': el('dryPath').value};
    const payload = {source_format: el('dryFormat').value, model: el('dryModel').value, requested_model: el('dryModel').value, stream: el('dryStream').value === 'true', headers: {}, body: Array.from(new TextEncoder().encode(JSON.stringify(reqBody))), metadata};
    const out = await readJSON(api.dryRun, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
    el('dryRunOut').textContent = JSON.stringify(out, null, 2);
    renderDryRunHints(reqBody);
    renderDryRunDecision(out);
    renderDryRunStageTrace(out.stage_trace || []);
    log('Dry-run completed.');
  }catch(err){ log(err.message, 'bad'); }
});
el('copyDryRunHintsBtn').addEventListener('click', async () => {
  try {
    const payload = currentDryRunHintPayload();
    if(!payload){ throw new Error('No dry-run payload to copy.'); }
    await navigator.clipboard.writeText(JSON.stringify(payload, null, 2));
    log('Dry-run hints copied.');
  } catch(err){ log(err.message || 'Clipboard copy failed.', 'bad'); }
});
function currentPrimaryMemberLabel(){
  const members = weightedRoutes();
  return members[0] ? (members[0].model || ((members[0].provider || '-') + '/' + (members[0].suffix || '-'))) : '';
}
function currentSecondaryMemberLabel(){
  const members = weightedRoutes();
  return members[1] ? (members[1].model || ((members[1].provider || '-') + '/' + (members[1].suffix || '-'))) : '';
}
function buildPoolPreviewPayload(action){
  if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
  const canaryPercent = Math.min(99, Math.max(1, Number(el('poolCanaryPercentInput')?.value || 10)));
  const shiftPercent = Math.min(100, Math.max(1, Number(el('poolShiftPercentInput')?.value || 50)));
  const primary = currentPrimaryMemberLabel();
  const secondary = el('poolCanarySecondaryInput')?.value?.trim() || currentSecondaryMemberLabel();
  switch(action){
    case 'pool-drain':
      return { key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'pool-drain', reason: 'pool-drain', preview_only: true };
    case 'pool-resume':
      return { key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'pool-resume', reason: 'pool-resume', preview_only: true };
    case 'canary-split':
      return { key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'canary-split', member: primary, secondary, primary_weight: 100 - canaryPercent, canary_weight: canaryPercent, reason: 'canary-split', preview_only: true };
    case 'rebalance-by-health':
      return { key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'rebalance-by-health', reason: 'rebalance-by-health', preview_only: true };
    case 'restore-default-weights':
      return { key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'restore-default-weights', reason: 'restore-default-weights', preview_only: true };
    case 'shift-provider-traffic':
      const provider = el('poolShiftProviderInput')?.value?.trim() || '';
      if(!provider){ throw new Error('Enter a provider to shift traffic to.'); }
      return { key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'shift-provider-traffic', secondary: provider, canary_weight: shiftPercent, reason: 'shift-provider-traffic', preview_only: true };
    default:
      throw new Error('Unsupported preview action.');
  }
}
function renderPoolPreview(result){
  const root = el('poolPreviewOut');
  if(!root) return;
  if(!result || typeof result !== 'object'){
    root.textContent = '{}';
    return;
  }
  const lines = [];
  lines.push('operation: ' + (result.operation || '-'));
  lines.push('target: ' + (result.target_member || '-'));
  lines.push('secondary: ' + (result.secondary || '-'));
  lines.push('reason: ' + (result.reason || '-'));
  lines.push('preview token: ' + (result.preview_token || '-'));
  lines.push('before: ' + (result.before_state || '-'));
  lines.push('after: ' + (result.after_state || '-'));
  if(Array.isArray(result.diff) && result.diff.length){
    lines.push('diff:');
    result.diff.forEach((item, index) => {
      const beforeWeight = item?.before?.weight ?? 0;
      const afterWeight = item?.after?.weight ?? 0;
      const beforeStatus = item?.before?.status || '-';
      const afterStatus = item?.after?.status || '-';
      const beforeHealth = item?.before?.health ?? 0;
      const afterHealth = item?.after?.health ?? 0;
      lines.push('  ' + (index + 1) + '. ' + (item.member || '-') + ' | weight ' + beforeWeight + ' -> ' + afterWeight + ' | status ' + beforeStatus + ' -> ' + afterStatus + ' | health ' + beforeHealth + ' -> ' + afterHealth);
    });
  }
  root.textContent = lines.join('\n');
  state.latestPoolPreviewToken = result.preview_token || '';
  if(el('poolPreviewTokenInput')) el('poolPreviewTokenInput').value = state.latestPoolPreviewToken;
}

async function previewThenApply(payload){
  const previewPayload = {...payload, preview_only: true};
  const preview = await readJSON(api.routeMemberPreview, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(previewPayload)});
  renderPoolPreview(preview);
  const token = preview?.preview_token || '';
  if(!token){
    throw new Error('Preview did not return a token.');
  }
  const applyPayload = {...payload, preview_token: token};
  delete applyPayload.preview_only;
  return readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(applyPayload)});
}

async function runPoolPreview(action){
  const payload = buildPoolPreviewPayload(action);
  const out = await readJSON(api.routeMemberPreview, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
  renderPoolPreview(out);
  log('Pool operation preview ready for ' + action + '.');
}
el('poolPreviewCanaryBtn').addEventListener('click', async () => {
  try { await runPoolPreview('canary-split'); } catch(err){ log(err.message || 'Canary preview failed.', 'bad'); }
});
el('poolPreviewShiftBtn').addEventListener('click', async () => {
  try { await runPoolPreview('shift-provider-traffic'); } catch(err){ log(err.message || 'Shift preview failed.', 'bad'); }
});
el('poolPreviewDrainBtn').addEventListener('click', async () => {
  try { await runPoolPreview('pool-drain'); } catch(err){ log(err.message || 'Drain preview failed.', 'bad'); }
});
el('poolPreviewResumeBtn').addEventListener('click', async () => {
  try { await runPoolPreview('pool-resume'); } catch(err){ log(err.message || 'Resume preview failed.', 'bad'); }
});
el('poolPreviewRebalanceBtn').addEventListener('click', async () => {
  try { await runPoolPreview('rebalance-by-health'); } catch(err){ log(err.message || 'Rebalance preview failed.', 'bad'); }
});
el('poolPreviewRestoreBtn').addEventListener('click', async () => {
  try { await runPoolPreview('restore-default-weights'); } catch(err){ log(err.message || 'Restore preview failed.', 'bad'); }
});
el('poolDrainBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'pool-drain', reason: 'pool-drain' });
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Route pool drained.');
  } catch(err){ log(err.message || 'Pool drain failed.', 'bad'); }
});
el('poolResumeBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'pool-resume', reason: 'pool-resume' });
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Route pool resumed.');
  } catch(err){ log(err.message || 'Pool resume failed.', 'bad'); }
});
el('poolCanaryBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    const members = weightedRoutes();
    if(members.length < 2){ throw new Error('Canary split needs at least two pool members.'); }
    const primary = members[0].model || ((members[0].provider || '-') + '/' + (members[0].suffix || '-'));
    const secondary = el('poolCanarySecondaryInput').value.trim() || (members[1].model || ((members[1].provider || '-') + '/' + (members[1].suffix || '-')));
    const canaryPercent = Math.min(99, Math.max(1, Number(el('poolCanaryPercentInput').value || 10)));
    await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'canary-split', member: primary, secondary, primary_weight: 100 - canaryPercent, canary_weight: canaryPercent, reason: 'canary-split' });
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Canary split applied: ' + primary + ' / ' + secondary + '.');
  } catch(err){ log(err.message || 'Canary split failed.', 'bad'); }
});
el('poolRebalanceBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'rebalance-by-health', reason: 'rebalance-by-health' });
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Route pool rebalanced by health.');
  } catch(err){ log(err.message || 'Health rebalance failed.', 'bad'); }
});
el('poolRestoreBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'restore-default-weights', reason: 'restore-default-weights' });
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Route pool weights restored.');
  } catch(err){ log(err.message || 'Restore weights failed.', 'bad'); }
});
el('poolShiftProviderBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    const provider = el('poolShiftProviderInput').value.trim();
    if(!provider){ throw new Error('Enter a provider to shift traffic to.'); }
    const percent = Math.min(100, Math.max(1, Number(el('poolShiftPercentInput').value || 50)));
    await previewThenApply({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'shift-provider-traffic', secondary: provider, canary_weight: percent, reason: 'shift-provider-traffic' });
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Shifted route traffic toward provider ' + provider + '.');
  } catch(err){ log(err.message || 'Shift provider failed.', 'bad'); }
});

buildTabbedLayout();
initManagementAccess();
</script>
</body>
</html>`
}
