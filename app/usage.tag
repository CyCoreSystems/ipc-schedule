<usage>

<h1 id="usage">Usage</h1>
<h2 id="general-use">General Use</h2>
<ul>
<li>GET <code>/</code> Print Instructions (this page)</li>
<li>GET <code>/target/:groupID</code> Print the current target for the given group ID; an optional <code>date</code> parameter may be passed to resolve the schedule for that date instead of now.</li>
</ul>
<h2 id="groups">Groups</h2>
<p>A <code>group</code> has the data structure:</p>
<div class="sourceCode"><pre class="sourceCode json"><code class="sourceCode json">            <span class="fu">{</span>
                <span class="dt">&quot;id&quot;</span><span class="fu">:</span> <span class="st">&quot;ID of group&quot;</span><span class="fu">,</span>
                <span class="dt">&quot;name&quot;</span><span class="fu">:</span> <span class="st">&quot;name/label of group&quot;</span><span class="fu">,</span>
                <span class="dt">&quot;timezone&quot;</span><span class="fu">:</span> <span class="st">&quot;Time zone, of the form US/Eastern or America/New York&quot;</span><span class="fu">,</span>
            <span class="fu">}</span></code></pre></div>
<ul>
<li><strong>GET</strong> <code>/group/:groupID</code> Print the group identified by groupID</li>
<li><strong>POST</strong> <code>/group</code> Add a group.</li>
</ul>
<h2 id="import">Import</h2>
<p>There are two types of CSV import: &quot;days&quot; and &quot;dates&quot;. &quot;days&quot; imports a default schedule, based<br />
on the provided generic days of the week. &quot;dates&quot; imports schedules for specific dates. If there<br />
exists a &quot;dates&quot; schedule for any given time, it is used in preference to the &quot;days&quot; schedule.</p>
<p>A &quot;days&quot; schedule is a CSV file with no field headers and columns of the form:</p>
<pre><code>   &quot;Group ID&quot;,&quot;Day of the Week&quot;,&quot;Start Time (HH:MM)&quot;,&quot;Stop Time (HH:MM)&quot;,&quot;Target phone number&quot;</code></pre>
<p><em>(Day of the Week can be one- or three-letter abbreviations or the full weekday name: 'M', 'Mon', 'Monday')</em></p>
<p>A &quot;dates&quot; schedule is a CSV file with no field headers and columns of the form:</p>
<pre><code>   &quot;Group ID&quot;,&quot;Date (YYYY-MM-DD)&quot;,&quot;Start Time (HH:MM)&quot;,&quot;Stop Time (HH:MM)&quot;,&quot;Target phone number&quot;</code></pre>
<ul>
<li><strong>POST</strong> <code>/sched/import/days</code> Add a days (generic weekly) schedule.</li>
<li><strong>POST</strong> <code>/sched/import/dates</code> Add a dates (specific dates) schedule.</li>
</ul>

</usage>
