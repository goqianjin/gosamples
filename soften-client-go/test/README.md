


## 测试用例

<table>
    <!-- header -->
    <tr>
        <th>类别</th>
        <th>用例</th>
        <th>场景</th>
        <th>备注</th>
    </tr>
    <!-- body -->
    <!-- produce cases -->
    <tr align="left">
        <th rowspan="4">Produce</th><th >send 1msg</th><th>单个消息</th><th>Passed</th>
    </tr>
    <tr align="left"><th>send 1msg (Ready) + send 1msg (route to L1)</th><th></th><th></th></tr>
    <tr align="left"><th>sendAsync 1msg (Ready) </th><th></th><th></th></tr>
    <tr align="left"><th>sendAsync 1msg (Ready) + send 1msg (route to L1) </th><th></th><th></th></tr>
    <!-- listen cases -->
    <tr align="left">
        <th rowspan="7">Listen</th></th><th>listen 1msg (Ready)</th><th></th><th>send msg before</th>
    </tr>
    <tr align="left"><th>listen 1msg (Pending)</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>listen 1msg (Blocking)</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>listen 1msg (Retrying)</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>listen 4msg from ALL status</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>listen 2msg from main(Ready) and L1</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>listen 3msg from main(Ready) and L1&B1</th><th></th><th>send msg before</th></tr>
    <!-- before check cases -->
    <tr align="left">
        <th rowspan="10">Before-Check</th></th><th>check goto Done 1msg</th><th></th><th>send msg before</th>
    </tr>
    <tr align="left"><th>check goto Discard 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Dead 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Pending 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Blocking 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Retrying 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Upgrade 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Degrade 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Reroute 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto ALL goto-actions 9msg</th><th></th><th>send msg before</th></tr>
    <!-- before check cases -->
    <tr align="left">
        <th rowspan="9">Handle</th></th><th>handle goto Done 1msg</th><th></th><th>send msg before</th>
    </tr>
    <tr align="left"><th>handle goto Discard 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto Dead 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto Pending 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto Blocking 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto Retrying 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto Upgrade 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto Degrade 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>handle goto ALL goto-actions 8msg</th><th></th><th>send msg before</th></tr>
    <!-- after check cases -->
    <tr align="left">
        <th rowspan="10">After-Check</th></th>
        <th>check goto Done 1msg</th><th></th><th>send msg before</th>
    </tr>
    <tr align="left"><th>check goto Discard 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Dead 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Pending 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Blocking 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Retrying 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Upgrade 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Degrade 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto Reroute 1msg</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>check goto ALL goto-actions 9msg</th><th></th><th>send msg before</th></tr>
    <!-- Status balance cases -->
    <tr align="left">
        <th rowspan="4">Status-Balance</th></th>
        <th>balance status messages: 50(Ready) + 30(Retrying)</th><th></th><th>send msg before</th>
    </tr>
    <tr align="left"><th>balance status messages: 50(Ready) + 15(Pending)</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>balance status messages: 50(Ready) + 5(Blocking)</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>balance status messages: ALL 50(Ready) + <br>30(Retrying) + 15(Pending) + 5(Blocking)</th><th></th><th>send msg before</th></tr>
    <!-- Level balance cases -->
    <tr align="left">
        <th rowspan="3">Level-Balance</th></th>
        <th>balance level messages: 50(L1) + 30(B1)</th><th></th><th>send msg before</th>
    </tr>
    <tr align="left"><th>balance level messages: 50(L1) + 100(L2)</th><th></th><th>send msg before</th></tr>
    <tr align="left"><th>balance status messages: ALL 50(L1) +<br> 100(L2) +30(B1)</th><th></th><th>send msg before</th></tr>
    <!-- Equality-Nursed cases -->
    <tr align="left">
        <th rowspan="5">Equality/Isolation</th></th>
        <th>1000(Normal) + 10000(Radical); <br>rate=100qps; <br>goto Retrying if exceeds the rate; <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th>
    </tr>
    <tr align="left"><th>1000(Normal) + 10000(Radical); <br>rate=100qps; <br>goto Retrying if exceeds the rate; <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th></tr>
    <tr align="left"><th>1000(Normal) + 10000(Radical); <br>rate=100qps; <br>goto Retrying if exceeds the rate; <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th></tr>
    <tr align="left"><th>1000(Normal) + 10000(Radical); <br>quota=5000; <br>goto Blocking if exceeds the rate; <br>Increase quota to 10000 after backlog is empty; <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th></tr>
    <tr align="left"><th>1000(Normal) + 10000(Radical); <br>quota=5000; <br>goto Blocking if exceeds the rate; <br>goto DLQ if exceeds max retries; <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th></tr>
    <!-- Overall cases -->
    <tr align="left">
        <th rowspan="2">Overall</th></th>
        <th>1000(Normal) + 10000(Radical); <br>rate=100qps; <br>goto Retrying if exceeds the rate; <br>Degrade to B1 If retries meets 5; <br> <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th>
    </tr>
    <tr align="left"><th>1000(Normal) + 10000(Radical) + 10000(Radical-2); <br>rate=100qps; <br>goto Retrying if exceeds the rate; <br>Degrade Radical to B1 If retries meets 5; <br>Upgrade Radical-2 to L2 if exceeds the rate; <br>Each normal msg is followed by 10 radical ones;</th><th></th><th>send msg before; <br>rate is used-based</th></tr>



</table>

