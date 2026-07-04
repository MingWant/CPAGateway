package main

func gatewayUIHTML() string {
	return `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta http-equiv="Content-Security-Policy" content="default-src 'none'; script-src 'nonce-gateway-ui'; style-src 'unsafe-inline'; connect-src 'self'; img-src data:; base-uri 'none'; form-action 'none'">
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
</style>
</head>
<body>
<div class="wrap">
  <section class="hero">
    <h1>Gateway Manager</h1>
    <p>Manage per-key gateway routing, model rewrite rules, daily limits, minute rate limits, and dry-run checks for CPA top-level API keys.</p>
    <div class="card" style="padding:12px 16px;border-radius:18px;position:sticky;top:12px;z-index:2"><strong>Current Context</strong><div id="contextBreadcrumb" class="hint">window 1h | no policy selected</div></div>
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
        <div class="actions"><button class="warn" id="addRouteToModelRuleBtn">Add Route Rule</button><button class="warn" id="addFallbackRuleBtn">Add Fallback Rule</button><button class="warn" id="addDenyRuleBtn">Add Deny Rule</button><button class="warn" id="buildFallbackChainBtn">Build Fallback Chain Rule</button></div><div class="card" style="padding:12px 0 0;border:0;box-shadow:none;background:none"><h3>Templates</h3><div class="grid2"><label>Template Search<input id="templateSearchInput" placeholder="name or description"></label><label>Template Category<select id="templateCategoryFilter"><option value="">all</option><option value="routing">routing</option><option value="fallback">fallback</option><option value="security">security</option><option value="custom">custom</option></select></label></div><div class="grid3"><label>Scenario<select id="templateScenarioFilter"><option value="">all</option><option value="model-migration">model-migration</option><option value="traffic-split">traffic-split</option><option value="cost-control">cost-control</option><option value="shadow-release">shadow-release</option><option value="provider-guardrail">provider-guardrail</option></select></label><label>Maturity<select id="templateMaturityFilter"><option value="">all</option><option value="stable">stable</option><option value="beta">beta</option><option value="experimental">experimental</option></select></label><label>Tag<input id="templateTagFilter" placeholder="routing"></label></div><div class="actions"><button id="saveTemplateBtn">Save Current Rule As Template</button><button class="alt" id="applyTemplateFilterBtn">Filter Templates</button><button class="alt" id="exportTemplatesBtn">Export Templates</button></div><label>Template Import JSON<textarea id="templateImportBox" placeholder="{"items":[...]}"></textarea></label><div class="actions"><button class="warn" id="importTemplatesBtn">Import Templates</button></div><div class="list" id="templateList"></div></div>
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
              <label>Any-Of JSON<textarea id="ruleAnyOfInput" placeholder="[{"models":["gpt-5.5"]},{"paths":["/v1/responses"]}]"></textarea></label>
              <label>All-Of JSON<textarea id="ruleAllOfInput" placeholder="[{"query":{"mode":"strict"}},{"body_contains":{"service_tier":"priority"}}]"></textarea></label>
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
const el = id => document.getElementById(id);
const initialURLState = readURLState();
let state = { keys: [], usage: [], audit: [], auditSummary: { total_by_decision: {}, total_by_reason: {}, total_by_rule: {}, total_by_policy: {}, total_by_model: {} }, templates: [], policies: { key_policies: [], default_policy: {} }, selectedKeyId: initialURLState.keyId || '', selectedRuleId: initialURLState.ruleId || '', focusedPreviewLabel: initialURLState.previewLabel || '', focusedPreviewType: initialURLState.previewType || '', latestPoolPreviewToken: '', weightedRoutes: [], failoverHops: [], memberHitCounts: {}, ruleHitCounts: {}, stageHitCounts: {}, memberHitCountsLast5m: {}, ruleHitCountsLast5m: {}, stageHitCountsLast5m: {}, memberHitCountsLastHour: {}, ruleHitCountsLastHour: {}, stageHitCountsLastHour: {}, memberHitCountsLast24h: {}, ruleHitCountsLast24h: {}, stageHitCountsLast24h: {}, statsWindow: initialURLState.statsWindow || '1h', collapsedTopology: Object.keys(initialURLState.collapsedTopology || {}).length ? initialURLState.collapsedTopology : loadCollapsedTopology() };
function log(msg, kind='ok'){ el('logBox').textContent = '[' + new Date().toLocaleTimeString() + '] ' + msg; el('logBox').style.color = kind === 'bad' ? '#ffb4ab' : '#9ff3e8'; }
async function readJSON(url, init){ const res = await fetch(url, init); const body = await res.json(); if(!res.ok) throw new Error(body.error || body.message || JSON.stringify(body)); return body; }
function selectedPolicy(){ return (state.policies.key_policies || []).find(item => item.key_id === state.selectedKeyId) || null; }
function renderKeys(){
  const root = el('keysList'); root.innerHTML='';
  const items = state.keys || [];
  if(!items.length){ root.innerHTML = '<div class="item">No key-specific policies yet.</div>'; return; }
  items.forEach(key => {
    const node = document.createElement('div'); node.className='item' + (state.selectedKeyId === key.key_id ? ' active' : '');
    node.innerHTML = '<strong>' + (key.display_name || key.key_id) + '</strong><span class="pill">' + (key.masked_key || '') + '</span><div class="hint">key_id: ' + key.key_id + ' | enabled: ' + key.enabled + '</div>';
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
function renderSubnodeList(type, items, options){
  const list = Array.isArray(items) ? items : [];
  const emptyText = options?.emptyText || 'No items.';
  if(!list.length){
    return emptyText ? '<div class="hint">' + emptyText + '</div>' : '';
  }
  return '<div class="subnode-list">' + list.map((item, index) => {
    const label = previewChipLabel(type, item);
    const focused = previewMatches(type, label) ? ' focused' : '';
    const summary = options?.summary ? options.summary(item, index) : '';
    const meta = options?.meta ? options.meta(item, index) : '';
    return '<div class="subnode ' + (options?.kind || '') + focused + '" data-preview-focus="' + type + '" data-preview-label="' + label.replace(/"/g, '&quot;') + '"><div><strong>' + label + '</strong>' + (summary ? '<span class="hint">' + summary + '</span>' : '') + '</div><div class="hint">' + (meta || '') + '</div></div>';
  }).join('') + '</div>';
}
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
  root.innerHTML = '';
  if(!payload || typeof payload !== 'object'){ root.innerHTML = '<div class="stat"><span>No dry-run hints yet.</span><b>0</b></div>'; return; }
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
    const node = document.createElement('div');
    node.className = 'stat';
    node.innerHTML = '<span>' + label + '</span><b>' + value + '</b>';
    root.appendChild(node);
  });
  if(!root.innerHTML){ root.innerHTML = '<div class="stat"><span>No dry-run hints yet.</span><b>0</b></div>'; }
}

function renderDryRunDecision(result){
  const root = el('dryRunDecisionStats');
  if(!root) return;
  root.innerHTML = '';
  if(!result || typeof result !== 'object'){
    root.innerHTML = '<div class="stat"><span>No dry-run decision yet.</span><b>0</b></div>';
    return;
  }
  [
    ['decision', result.decision || 'pass'],
    ['reason', result.reason || '-'],
    ['final_model', result.final_model || '-'],
    ['rule_id', result.rule_id || '-']
  ].forEach(([label, value]) => {
    const node = document.createElement('div');
    node.className = 'stat';
    node.innerHTML = '<span>' + label + '</span><b>' + value + '</b>';
    root.appendChild(node);
  });
}

function renderDryRunStageTrace(items){
  const root = el('dryRunStageTrace');
  if(!root) return;
  root.innerHTML = '';
  const rows = Array.isArray(items) ? items : [];
  if(!rows.length){
    root.innerHTML = '<div class="stat"><span>No stage trace yet.</span><b>0</b></div>';
    return;
  }
  rows.forEach((item, index) => {
    const node = document.createElement('div');
    node.className = 'item';
    const failoverReasons = Array.isArray(item.failover_reasons) && item.failover_reasons.length ? item.failover_reasons.join(', ') : '-';
    const mirrors = Array.isArray(item.mirror_models) && item.mirror_models.length ? item.mirror_models.join(', ') : '-';
    const failoverChain = Array.isArray(item.failover_chain) && item.failover_chain.length ? item.failover_chain.join(' -> ') : '-';
    node.innerHTML = '<strong>' + (index + 1) + '. ' + (item.stage || '-') + '</strong>' +
      '<div class="hint">mode ' + (item.mode || '-') + ' | matched ' + (item.matched_count ?? (item.matched_rules || []).length || 0) + ' | decision ' + (item.decision || 'pass') + '</div>' +
      '<div class="hint">rules ' + (((item.matched_rules || []).length ? item.matched_rules.join(', ') : '-')) + '</div>' +
      '<div class="hint">final model ' + (item.final_model || '-') + ' | route target ' + (item.route_target || '-') + '</div>' +
      '<div class="hint">route pool ' + (item.route_pool || '-') + ' | fallback target ' + (item.fallback_target || '-') + '</div>' +
      '<div class="hint">mirrors ' + mirrors + '</div>' +
      '<div class="hint">failover chain ' + failoverChain + '</div>' +
      '<div class="hint">failover reasons ' + failoverReasons + '</div>' +
      '<div class="hint">reason ' + (item.reason || '-') + '</div>';
    root.appendChild(node);
  });
}

function renderCurrentNodeDetail(){
  const root = el('currentNodeDetail');
  if(!root) return;
  root.innerHTML = '';
  const policy = selectedPolicy();
  const rule = (policy?.rules || []).find(item => item.id === state.selectedRuleId);
  if(!policy){ root.innerHTML = '<div class="item">No policy selected.</div>'; return; }
  const windowHits = currentWindowHitMaps();
  const policyCard = document.createElement('div');
  policyCard.className = 'item';
  policyCard.innerHTML = '<strong>Policy</strong><div class="hint">' + (policy.display_name || policy.key_id) + '</div><div class="hint">enabled ' + (policy.enabled !== false) + ' | rules ' + ((policy.rules || []).length) + ' | window ' + windowHits.label + '</div>';
  root.appendChild(policyCard);
  const governanceCard = document.createElement('div');
  governanceCard.className = 'item';
  governanceCard.innerHTML = '<strong>Governance</strong><div class="hint">requests/day ' + (policy.limits?.requests_per_day || 0) + ' | requests/min ' + (policy.limits?.requests_per_minute || 0) + ' | inflight ' + (policy.limits?.max_inflight || 0) + '</div><div class="hint">pre-check ' + (policy.stage_policy?.['pre-check']?.mode || 'first-match') + ' | route ' + (policy.stage_policy?.route?.mode || 'first-match') + ' | mirror ' + (policy.stage_policy?.mirror?.mode || 'continue-all') + '</div>';
  root.appendChild(governanceCard);
  const quickCard = document.createElement('div');
  quickCard.className = 'item';
  quickCard.innerHTML = '<strong>Quick Actions</strong><div class="actions"><button class="alt" id="focusPolicyBtn">Focus Policy</button><button class="alt" id="focusRuleBtn">Focus Rule</button><button class="alt" id="openDryRunBtn">Open Dry-Run</button><button class="alt" id="resetDryRunBtn">Reset Dry-Run</button></div><div class="actions"><button class="alt" id="copyDryRunRequestBtn">Copy Request</button><button class="alt" id="copyRuleSummaryBtn">Copy Summary</button><button class="alt" id="copyRuleIdBtn">Copy Rule ID</button><button class="alt" id="copyRoutePoolBtn">Copy Pool</button><button class="alt" id="copyFailoverChainBtn">Copy Failover</button></div>';
  root.appendChild(quickCard);
  setTimeout(() => {
    const focusPolicyBtn = el('focusPolicyBtn');
    if(focusPolicyBtn) focusPolicyBtn.onclick = () => { renderKeys(); renderRules(); renderRouteGraph(); log('Focused current policy in topology.'); };
    const focusRuleBtn = el('focusRuleBtn');
    if(focusRuleBtn) focusRuleBtn.onclick = () => { if(rule){ hydrateRule(rule); renderRules(); renderRouteGraph(); log('Focused current rule.'); } };
    const openDryRunBtn = el('openDryRunBtn');
    if(openDryRunBtn) openDryRunBtn.onclick = () => {
      el('dryKey').value = policy.match_api_key || '';
      if(rule?.match?.models?.[0]) el('dryModel').value = rule.match.models[0];
      if(rule?.match?.paths?.[0]) el('dryPath').value = rule.match.paths[0];
      applyDryRunPayloadFromRule(policy, rule);
      log('Prepared dry-run form from current node.');
    };
  }, 0);
  if(rule){
    const ruleCard = document.createElement('div');
    ruleCard.className = 'item';
    ruleCard.innerHTML = '<strong>Rule</strong><div class="hint">' + (rule.id || '-') + ' | stage ' + (rule.stage || 'pre-check') + '</div><div class="hint">priority ' + (rule.priority || 0) + ' | total hits ' + (state.ruleHitCounts[rule.id || ''] || 0) + ' | current window ' + (windowHits.rule[rule.id || ''] || 0) + '</div>';
    root.appendChild(ruleCard);
    const actionCard = document.createElement('div');
    actionCard.className = 'item';
    actionCard.innerHTML = '<strong>Route</strong><div class="hint">pool ' + (rule.actions?.route_pool?.name || '-') + ' | route_to ' + (rule.actions?.route_to_model || '-') + '</div><div class="hint">failover hops ' + ((rule.actions?.failover_hops || []).length) + ' | mirrors ' + ((rule.actions?.mirror_models || []).length) + '</div>';
    root.appendChild(actionCard);
    const poolMembers = rule.actions?.route_pool?.members || rule.actions?.weighted_routes || [];
    if(poolMembers.length){
      const poolCard = document.createElement('div');
      poolCard.className = 'item';
      const previewItems = poolMembers.slice(0, 3).map(member => { const label = previewChipLabel('pool-member', member); const active = previewMatches('pool-member', label) ? ' style="outline:2px solid #a16207;outline-offset:2px"' : ''; return '<span class="chip"' + active + ' data-preview-focus="pool-member" data-preview-label="' + label.replace(/"/g, '&quot;') + '">' + label + '</span>'; }).join(' ');
      poolCard.innerHTML = '<strong>Pool Members</strong><div class="chips">' + previewItems + (poolMembers.length > 3 ? '<span class="chip">...</span>' : '') + '</div><div class="hint">Click a member chip to refocus the topology route branch for the selected rule.</div>';
      poolCard.title = 'Use the route branch in topology to inspect pool members for the selected rule.';
      bindPreviewFocus(poolCard, policy, rule);
      root.appendChild(poolCard);
    }
    if((rule.actions?.failover_hops || []).length){
      const hopCard = document.createElement('div');
      hopCard.className = 'item';
      const previewItems = (rule.actions.failover_hops || []).slice(0, 3).map(hop => { const label = previewChipLabel('failover-hop', hop); const active = previewMatches('failover-hop', label) ? ' style="outline:2px solid #a16207;outline-offset:2px"' : ''; return '<span class="chip"' + active + ' data-preview-focus="failover-hop" data-preview-label="' + label.replace(/"/g, '&quot;') + '">' + label + '</span>'; }).join(' ');
      hopCard.innerHTML = '<strong>Failover Hops</strong><div class="chips">' + previewItems + ((rule.actions.failover_hops || []).length > 3 ? '<span class="chip">...</span>' : '') + '</div><div class="hint">Click a hop chip to localize the exact failover target in topology order.</div>';
      hopCard.title = 'Use the failover branch in topology to inspect hop order for the selected rule.';
      bindPreviewFocus(hopCard, policy, rule);
      root.appendChild(hopCard);
    }
    if((rule.actions?.mirror_models || []).length){
      const mirrorCard = document.createElement('div');
      mirrorCard.className = 'item';
      const previewItems = (rule.actions.mirror_models || []).slice(0, 4).map(model => { const label = previewChipLabel('mirror-target', { model }); const active = previewMatches('mirror-target', label) ? ' style="outline:2px solid #a16207;outline-offset:2px"' : ''; return '<span class="chip"' + active + ' data-preview-focus="mirror-target" data-preview-label="' + label.replace(/"/g, '&quot;') + '">' + label + '</span>'; }).join(' ');
      mirrorCard.innerHTML = '<strong>Mirror Targets</strong><div class="chips">' + previewItems + ((rule.actions.mirror_models || []).length > 4 ? '<span class="chip">...</span>' : '') + '</div><div class="hint">Mirror targets are shown as shadow traffic branches in topology.</div>';
      bindPreviewFocus(mirrorCard, policy, rule);
      root.appendChild(mirrorCard);
    }
  } else {
    const info = document.createElement('div');
    info.className = 'item';
    info.innerHTML = '<strong>Rule</strong><div class="hint">No rule selected in this policy.</div>';
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
  const root = el('usageStats'); root.innerHTML='';
  const rows = state.usage || [];
  const windowHits = currentWindowHitMaps();
  renderStatsWindowButtons();
  renderContextBreadcrumb();
  renderCurrentNodeDetail();
  if(!rows.length){ root.innerHTML = '<div class="stat"><span>No usage data yet.</span><b>0</b></div>'; renderGlobalTopology(); return; }
  rows.forEach(row => {
    const node = document.createElement('div'); node.className='stat';
    node.innerHTML = '<span>' + (row.display_name || row.key_id) + '</span><b>' + row.requests_today + '</b><div class="hint">minute ' + row.requests_minute + ' | inflight ' + row.inflight + ' | window ' + windowHits.label + '</div>';
    root.appendChild(node);
  });
  renderGlobalTopology();
}
function renderAudit(){
  const root = el('auditList'); root.innerHTML='';
  const rows = state.audit || [];
  if(!rows.length){ root.innerHTML = '<div class="item">No audit records yet.</div>'; return; }
  rows.forEach(row => {
    const node = document.createElement('div'); node.className='item';
    node.innerHTML = '<strong>' + (row.decision || 'pass') + '</strong><div class="hint">' + (row.rule_id || '-') + ' | ' + (row.reason || '-') + '</div><div class="hint">' + (row.requested_model || '-') + ' -> ' + (row.final_model || '-') + '</div><div class="hint">event ' + (row.event_type || 'request') + ' | operator ' + (row.operator_action || '-') + ' | member ' + (row.target_member || '-') + '</div><div class="hint">before ' + (row.before_state || '-') + '</div><div class="hint">after ' + (row.after_state || '-') + '</div>';
    node.addEventListener('click', async () => { try { const detail = await readJSON(api.auditDetail + '?time=' + encodeURIComponent(row.time || '') + '&rule=' + encodeURIComponent(row.rule_id || '') + '&reason=' + encodeURIComponent(row.reason || '')); el('auditDetailOut').textContent = JSON.stringify(detail, null, 2); const stats = el('auditDetailStats'); stats.innerHTML = ''; ['decision','rule_id','reason','requested_model','final_model','event_type','operator_action','target_member','secondary','before_state','after_state'].forEach(key => { const value = detail[key]; if(!value) return; const nodeStat = document.createElement('div'); nodeStat.className='stat'; nodeStat.innerHTML = '<span>' + key + '</span><b>' + value + '</b>'; stats.appendChild(nodeStat); }); if(Array.isArray(detail.diff) && detail.diff.length){ const diffNode = document.createElement('div'); diffNode.className='stat'; diffNode.innerHTML = '<span>diff</span><b>' + detail.diff.length + ' changes</b>'; stats.appendChild(diffNode); } log('Audit detail loaded.'); } catch(err){ log(err.message, 'bad'); } });
    root.appendChild(node);
  });
}
function renderAuditSummary(){
  const root = el('auditSummaryStats'); root.innerHTML='';
  const groups = state.auditSummary?.total_by_decision || {};
  const byPolicy = state.auditSummary?.total_by_policy || {};
  const byModel = state.auditSummary?.total_by_model || {};
  if(!Object.keys(groups).length && !Object.keys(byPolicy).length && !Object.keys(byModel).length){ root.innerHTML = '<div class="stat"><span>No audit summary yet.</span><b>0</b></div>'; return; }
  Object.entries(groups).forEach(([key, value]) => {
    const node = document.createElement('div'); node.className='stat';
    node.innerHTML = '<span>decision: ' + key + '</span><b>' + value + '</b>';
    node.addEventListener('click', () => { el('auditDecisionFilter').value = key; refreshAll(); });
    root.appendChild(node);
  });
  Object.entries(byPolicy).slice(0, 3).forEach(([key, value]) => {
    const node = document.createElement('div'); node.className='stat';
    node.innerHTML = '<span>policy: ' + key + '</span><b>' + value + '</b>';
    root.appendChild(node);
  });
  Object.entries(byModel).slice(0, 3).forEach(([key, value]) => {
    const node = document.createElement('div'); node.className='stat';
    node.innerHTML = '<span>final model: ' + key + '</span><b>' + value + '</b>';
    root.appendChild(node);
  });
  const timelineRoot = el('auditTimelineStats');
  if(timelineRoot){
    timelineRoot.innerHTML = '';
    (state.auditSummary?.timeline || []).slice(-8).forEach(item => {
      const node = document.createElement('div'); node.className='stat';
      node.innerHTML = '<span>' + item.window + '</span><b>' + item.count + '</b>';
      timelineRoot.appendChild(node);
    });
  }
}
function renderTemplates(){
  const root = el('templateList'); root.innerHTML='';
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
  if(!items.length){ root.innerHTML = '<div class="item">No templates loaded.</div>'; return; }
  items.forEach(item => {
    const node = document.createElement('div'); node.className='item';
    node.innerHTML = '<strong>' + item.name + '</strong><div class="hint">' + item.description + '</div><div class="hint">category: ' + (item.category || 'custom') + ' | scenario: ' + (item.scenario || '-') + ' | maturity: ' + (item.maturity || '-') + '</div><div class="chips">' + ((item.tags || []).map(tag => '<span class="chip">' + tag + '</span>').join('') || '<span class="chip">no-tags</span>') + '</div><div class="hint">id: ' + item.id + '</div>';
    node.addEventListener('click', () => { const cloned = JSON.parse(JSON.stringify(item.rule)); if(!cloned.id){ cloned.id = 'template-rule-' + Date.now(); } appendRuleTemplate(cloned); hydrateRule(cloned); log('Template appended: ' + item.name); });
    const editBtn = document.createElement('button'); editBtn.className='alt'; editBtn.textContent='Edit'; editBtn.addEventListener('click', async (event) => { event.stopPropagation(); const name = prompt('Template name', item.name); if(name === null) return; const description = prompt('Template description', item.description || ''); if(description === null) return; const scenario = prompt('Template scenario', item.scenario || ''); if(scenario === null) return; const maturity = prompt('Template maturity', item.maturity || 'stable'); if(maturity === null) return; const tags = prompt('Template tags (comma separated)', (item.tags || []).join(',')); if(tags === null) return; try { await readJSON(api.templates + '?template_id=' + encodeURIComponent(item.id), {method:'PATCH', headers:{'Content-Type':'application/json'}, body: JSON.stringify({ ...item, name, description, scenario, maturity, tags: parseCSV(tags) })}); log('Template updated.'); await refreshAll(); } catch(err){ log(err.message, 'bad'); } });
    const cloneBtn = document.createElement('button'); cloneBtn.className='warn'; cloneBtn.textContent='Clone'; cloneBtn.addEventListener('click', async (event) => { event.stopPropagation(); try { await readJSON(api.templates + '/clone?template_id=' + encodeURIComponent(item.id), {method:'POST'}); log('Template cloned.'); await refreshAll(); } catch(err){ log(err.message, 'bad'); } });
    const useBtn = document.createElement('button'); useBtn.className='alt'; useBtn.textContent='Use'; useBtn.addEventListener('click', (event) => { event.stopPropagation(); appendRuleTemplate(item.rule); log('Template inserted into current policy: ' + item.name); });
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

function failoverHops(){ return Array.isArray(state.failoverHops) ? state.failoverHops : []; }
function setFailoverHops(items){ state.failoverHops = Array.isArray(items) ? items.slice() : []; renderFailoverHops(); }
function renderFailoverHops(){
  const root = el('failoverHopsList');
  if(!root) return;
  root.innerHTML = '';
  const items = failoverHops();
  if(!items.length){ root.innerHTML = '<div class="item">No failover hops configured.</div>'; return; }
  items.forEach((hop, index) => {
    const node = document.createElement('div'); node.className='tableitem';
    const label = hop.model || ((hop.provider || '-') + '/' + (hop.suffix || '-'));
    node.innerHTML = '<div><b>' + label + '</b><span>' + (hop.on_decision || 'reject') + ' | ' + (hop.reason || '-') + ' | enabled ' + String(hop.enabled !== false) + '</span></div><span>hop ' + (index + 1) + '</span>';
    const actions = document.createElement('div'); actions.className='actions';
    const editBtn = document.createElement('button'); editBtn.className='alt'; editBtn.textContent='Edit'; editBtn.addEventListener('click', () => { el('failoverHopModelInput').value = hop.model || ''; el('failoverHopProviderInput').value = hop.provider || ''; el('failoverHopSuffixInput').value = hop.suffix || ''; el('failoverHopReasonInput').value = hop.reason || ''; el('failoverHopDecisionInput').value = hop.on_decision || 'reject'; el('failoverHopEnabledInput').value = String(hop.enabled !== false); state.failoverHops.splice(index, 1); renderFailoverHops(); log('Failover hop loaded into editor.'); });
    const deleteBtn = document.createElement('button'); deleteBtn.className='danger'; deleteBtn.textContent='Delete'; deleteBtn.addEventListener('click', () => { state.failoverHops.splice(index, 1); renderFailoverHops(); log('Failover hop removed.'); });
    actions.append(editBtn, deleteBtn); node.appendChild(actions); root.appendChild(node);
  });
  el('ruleFailoverHopsInput').value = JSON.stringify(items, null, 2);
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
  root.innerHTML = '';
  root.className = 'topology-flow';
  const policies = state.policies?.key_policies || [];
  if(!policies.length){ root.innerHTML = '<div class="item">No policy topology yet.</div>'; return; }
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
    policyNode.innerHTML = '<strong>' + (policy.display_name || policy.key_id) + '</strong><div class="hint">rules ' + rules.length + ' | pools ' + pools + ' | failovers ' + failovers + ' | mirrors ' + mirrors + '</div><div class="' + statTone(totalRuleHitsHour) + '">hits total ' + totalRuleHits + ' | 5m ' + totalRuleHits5m + ' | 1h ' + totalRuleHitsHour + ' | 24h ' + totalRuleHits24h + ' | current ' + windowHits.label + '</div><div class="hint">stage modes: pre-check ' + (policy.stage_policy?.['pre-check']?.mode || 'first-match') + ' -> route ' + (policy.stage_policy?.route?.mode || 'first-match') + ' | route-stage hits ' + routeStageHits + '</div><div class="hint">' + (collapsed ? '[+] expand policy' : '[-] collapse policy') + '</div><div class="node-detail"><b>Policy Detail</b><span>Key id: ' + (policy.key_id || '-') + '</span><span>Display: ' + (policy.display_name || '-') + '</span><span>Limits: day ' + (policy.limits?.requests_per_day || 0) + ', min ' + (policy.limits?.requests_per_minute || 0) + ', inflight ' + (policy.limits?.max_inflight || 0) + '</span></div>';
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
        stageNode.innerHTML = '<strong>Stage: ' + stage + '</strong><div class="' + statTone(state.stageHitCountsLastHour[stage] || 0) + '">rules ' + stageRules.length + ' | hits ' + (state.stageHitCounts[stage] || 0) + ' | 5m ' + (state.stageHitCountsLast5m[stage] || 0) + ' | 1h ' + (state.stageHitCountsLastHour[stage] || 0) + ' | 24h ' + (state.stageHitCountsLast24h[stage] || 0) + '</div><div class="hint">' + (stageCollapsed ? '[+] expand stage' : '[-] collapse stage') + '</div><div class="node-detail"><b>Stage Detail</b><span>Current window hits: ' + (windowHits.stage[stage] || 0) + '</span><span>Stage mode: ' + ((policy.stage_policy?.[stage]?.mode) || (stage === 'mirror' || stage === 'post-audit' ? 'continue-all' : 'first-match')) + '</span><span>Actions</span><div class="actions"><button class="alt" data-topology-action="toggle-stage">Toggle Stage</button><button class="alt" data-topology-action="copy-stage-summary">Copy Summary</button></div></div>';
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
            const subnodes = renderSubnodeList('pool-member', poolMembers, { kind: 'route', emptyText: 'No route pool members.', summary: (member, index) => 'member ' + (index + 1) + ' | weight ' + (member.weight ?? 1) + ' | priority ' + (member.priority ?? 100), meta: member => 'status ' + (member.status || 'active') + ' | health ' + (member.health ?? 100) + ' | cap ' + (member.traffic_cap ?? 100) }) + renderSubnodeList('failover-hop', failoverItems, { kind: 'failover', emptyText: '', summary: (hop, index) => 'hop ' + (index + 1) + ' | decision ' + (hop.on_decision || 'reject'), meta: hop => (hop.reason || 'fallback') + ' | enabled ' + String(hop.enabled !== false) }) + renderSubnodeList('mirror-target', mirrorItems, { kind: 'mirror', emptyText: '', summary: (_item, index) => 'mirror ' + (index + 1), meta: () => 'shadow traffic target' });
            child.className = 'topology-rule' + (state.selectedRuleId === (rule.id || '') ? ' active' : '') + (previewFocused ? ' preview-focus' : '');
            child.innerHTML = '<div class="item"><strong>Rule: ' + (rule.id || '-') + '</strong><div class="' + statTone(state.ruleHitCountsLastHour[rule.id || ''] || 0) + '">stage ' + stage + ' | total ' + (state.ruleHitCounts[rule.id || ''] || 0) + ' | 5m ' + (state.ruleHitCountsLast5m[rule.id || ''] || 0) + ' | 1h ' + (state.ruleHitCountsLastHour[rule.id || ''] || 0) + ' | 24h ' + (state.ruleHitCountsLast24h[rule.id || ''] || 0) + '</div><div class="node-detail"><b>Rule Detail</b><span>Model: ' + ((rule.match?.models || [])[0] || '-') + '</span><span>Route pool: ' + (rule.actions?.route_pool?.name || '-') + '</span><span>Failover hops: ' + ((rule.actions?.failover_hops || []).length) + ' | Mirrors: ' + ((rule.actions?.mirror_models || []).length) + '</span>' + (previewLabel ? '<span>Preview focus: ' + previewLabel + '</span>' : '') + '<span>Targets</span>' + subnodes + '<span>Actions</span><div class="actions"><button class="alt" data-topology-action="focus-rule">Focus Rule</button><button class="alt" data-topology-action="open-rule-dry-run">Dry-Run</button><button class="alt" data-topology-action="copy-rule-summary-inline">Copy Summary</button></div></div></div>';
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
  root.innerHTML = '';
  const members = weightedRoutes();
  const usage = state.usage || [];
  const audit = state.audit || [];
  if(!members.length){ root.innerHTML = '<div class="item">No pool members configured.</div>'; return; }
  const latestUsage = usage[0] || null;
  members.forEach((member, index) => {
    const label = member.model || ((member.provider || '-') + '/' + (member.suffix || '-'));
    const enabled = member.enabled !== false;
    const windowHits = currentWindowHitMaps();
    const hitCount = state.memberHitCounts[String(label).toLowerCase()] || audit.filter(item => (item.final_model || '').toLowerCase() === String(label).toLowerCase()).length;
    const hitWindow = windowHits.member[String(label).toLowerCase()] || 0;
    const node = document.createElement('div');
    node.className = 'item';
    node.title = 'Pool member health comes from configured values plus recent hit counts in the selected window.';
    node.innerHTML = '<strong>' + label + '</strong><div class="hint">status ' + (member.status || 'active') + ' | enabled ' + enabled + ' | priority ' + (member.priority ?? 100) + '</div><div class="hint">health ' + (member.health ?? 100) + ' | traffic cap ' + (member.traffic_cap ?? 100) + ' | weight ' + (member.weight ?? 1) + '</div><div class="hint">hits ' + hitCount + ' | current ' + windowHits.label + ' ' + hitWindow + (latestUsage ? ' | inflight ' + latestUsage.inflight + ' | minute ' + latestUsage.requests_minute : '') + '</div><div class="hint">' + (member.reason || 'no member note') + '</div><div class="actions"><button class="alt" data-member-op="focus" data-member-index="' + index + '">Focus</button><button class="alt" data-member-op="active" data-member-index="' + index + '">Active</button><button class="warn" data-member-op="drain" data-member-index="' + index + '">Drain</button><button class="danger" data-member-op="offline" data-member-index="' + index + '">Offline</button><button class="alt" data-member-op="cap-down" data-member-index="' + index + '">Cap -10</button><button class="alt" data-member-op="cap-up" data-member-index="' + index + '">Cap +10</button></div><div class="actions"><button class="alt" data-member-op="weight-down" data-member-index="' + index + '">Weight -1</button><button class="alt" data-member-op="weight-up" data-member-index="' + index + '">Weight +1</button><button class="alt" data-member-op="priority-down" data-member-index="' + index + '">Priority -10</button><button class="alt" data-member-op="priority-up" data-member-index="' + index + '">Priority +10</button><button class="alt" data-member-op="health-down" data-member-index="' + index + '">Health -10</button><button class="alt" data-member-op="health-up" data-member-index="' + index + '">Health +10</button></div>';
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
        await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, member: label, operation: op, delta: 10, reason: op === 'drain' ? 'manual-drain' : (op === 'offline' ? 'manual-offline' : '') })});
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
  root.innerHTML = '';
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
      node.innerHTML = '<strong>' + (index + 1) + '. ' + block.title + '</strong><div class="' + badge + '">' + block.body + '</div><div class="hint">' + (block.meta || 'Request enters gateway rule and flows to the next node.') + '</div>' + (index < blocks.length - 1 ? '<div class="hint">|</div><div class="hint">v</div>' : '');
      root.appendChild(node);
    });
  } catch {
    root.innerHTML = '<div class="item">No route graph yet.</div>';
  }
}
function renderRules(){
  const root = el('rulesList'); root.innerHTML='';
  const rules = selectedPolicy()?.rules || [];
  if(!rules.length){ root.innerHTML = '<div class="item">No rules yet.</div>'; return; }
  rules.forEach(rule => {
    const node = document.createElement('div'); node.className='item' + (state.selectedRuleId === rule.id ? ' active' : '');
    const action = rule.actions?.route_to_model || (rule.actions?.route_pool?.name || '') || (rule.actions?.fallback_models || [])[0] || (rule.actions?.deny ? 'deny' : 'pass');
    const anyCount = (rule.match?.any_of || []).length;
    const allCount = (rule.match?.all_of || []).length;
    node.innerHTML = '<strong>' + (rule.id || 'rule') + '</strong><div class="hint">priority ' + (rule.priority || 0) + ' | on_match ' + (rule.on_match || 'stop') + '</div><div class="hint">action: ' + action + '</div><div class="hint">groups: any=' + anyCount + ' / all=' + allCount + '</div>';
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
  el('ruleFallbackModelInput').value = (rule.actions?.fallback_models || []).join(',');
  el('ruleFailoverChainInput').value = (rule.actions?.failover_chain || []).join(',');
  el('routePoolNameInput').value = rule.actions?.route_pool?.name || '';
  el('routePoolModeInput').value = rule.actions?.route_pool?.mode || 'weighted';
  el('routePoolAffinityInput').value = rule.actions?.route_pool?.provider_affinity || '';
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
    if(!state.selectedKeyId && state.keys.length){ state.selectedKeyId = state.keys[0].key_id; }
    syncURLState();
    if(state.selectedKeyId && !(state.keys || []).some(item => item.key_id === state.selectedKeyId)){ state.selectedKeyId = state.keys[0]?.key_id || ''; }
    el('policiesBox').value = JSON.stringify(state.policies, null, 2);
    renderKeys();
    renderUsage();
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
    log('Gateway data refreshed.');
  }catch(err){ log(err.message, 'bad'); }
}
el('refreshBtn').addEventListener('click', refreshAll);
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
el('addWeightedRouteBtn').addEventListener('click', () => { const model = el('weightedRouteModelInput').value.trim(); const provider = el('weightedRouteProviderInput').value.trim(); const suffix = el('weightedRouteSuffixInput').value.trim(); if(!model && !(provider && suffix)){ log('Enter a weighted route model or provider+suffix first.', 'bad'); return; } const weight = Number(el('weightedRouteWeightInput').value || 1) || 1; const priority = Number(el('weightedRoutePriorityInput').value || 100) || 100; const enabled = el('weightedRouteEnabledInput').value === 'true'; const status = el('weightedRouteStatusInput').value || 'active'; const reason = el('weightedRouteReasonInput').value.trim(); const health = Number(el('weightedRouteHealthInput').value || 100) || 100; const trafficCap = Number(el('weightedRouteTrafficCapInput').value || 100) || 100; state.weightedRoutes.push({ model, provider, suffix, weight, priority, enabled, status, reason, health, traffic_cap: trafficCap }); el('weightedRouteModelInput').value=''; el('weightedRouteProviderInput').value=''; el('weightedRouteSuffixInput').value=''; el('weightedRouteWeightInput').value='1'; el('weightedRoutePriorityInput').value='100'; el('weightedRouteEnabledInput').value='true'; el('weightedRouteStatusInput').value='active'; el('weightedRouteReasonInput').value=''; el('weightedRouteHealthInput').value='100'; el('weightedRouteTrafficCapInput').value='100'; renderWeightedRoutes(); log('Weighted route added.'); });
el('clearWeightedRoutesBtn').addEventListener('click', () => { state.weightedRoutes = []; renderWeightedRoutes(); log('Weighted routes cleared.'); });
el('sortWeightedRoutesBtn').addEventListener('click', () => { state.weightedRoutes = weightedRoutes().slice().sort((a, b) => a.model.localeCompare(b.model)); renderWeightedRoutes(); log('Weighted routes sorted.'); });
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
  if(weightedRoutes().length){ rule.actions.weighted_routes = weightedRoutes(); rule.actions.route_pool = { name: routePoolName, mode: routePoolMode || 'weighted', provider_affinity: routePoolAffinity, members: weightedRoutes() }; }
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

function poolApplyGuardPayload(payload){
  const safeApply = el('poolSafeApplyInput')?.value === 'true';
  if(safeApply && !state.latestPoolPreviewToken){
    throw new Error('Safe Apply is enabled. Run preview first.');
  }
  if(safeApply){
    payload.preview_token = state.latestPoolPreviewToken;
  }
  return payload;
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
    await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(poolApplyGuardPayload({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'pool-drain', reason: 'pool-drain' }))});
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Route pool drained.');
  } catch(err){ log(err.message || 'Pool drain failed.', 'bad'); }
});
el('poolResumeBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(poolApplyGuardPayload({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'pool-resume', reason: 'pool-resume' }))});
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
    await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(poolApplyGuardPayload({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'canary-split', member: primary, secondary, primary_weight: 100 - canaryPercent, canary_weight: canaryPercent, reason: 'canary-split' }))});
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Canary split applied: ' + primary + ' / ' + secondary + '.');
  } catch(err){ log(err.message || 'Canary split failed.', 'bad'); }
});
el('poolRebalanceBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(poolApplyGuardPayload({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'rebalance-by-health', reason: 'rebalance-by-health' }))});
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Route pool rebalanced by health.');
  } catch(err){ log(err.message || 'Health rebalance failed.', 'bad'); }
});
el('poolRestoreBtn').addEventListener('click', async () => {
  try {
    if(!state.selectedKeyId || !state.selectedRuleId){ throw new Error('Select a route rule first.'); }
    await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(poolApplyGuardPayload({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'restore-default-weights', reason: 'restore-default-weights' }))});
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
    await readJSON(api.routeMemberOp, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(poolApplyGuardPayload({ key_id: state.selectedKeyId, rule_id: state.selectedRuleId, operation: 'shift-provider-traffic', secondary: provider, canary_weight: percent, reason: 'shift-provider-traffic' }))});
    await refreshAll();
    const rule = (selectedPolicy()?.rules || []).find(item => item.id === state.selectedRuleId);
    if(rule){ hydrateRule(rule); }
    log('Shifted route traffic toward provider ' + provider + '.');
  } catch(err){ log(err.message || 'Shift provider failed.', 'bad'); }
});

refreshAll();
</script>
</body>
</html>`
}
