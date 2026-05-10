package ui

const HomeTemplate = `
<window id="home" title="Freecode" center="true" padding="2">
  <vbox gap="1">
    <text id="banner" value="${banner}" bold="true" />
    <spacer height="2" />
    <text id="input-label" value="${inputLabel}" />
    <spacer flex="1" />
    <text id="hints" value="Ctrl+P: Command Palette | Ctrl+B: Toggle Sidebar | Ctrl+H: Home | Ctrl+Q: Quit" color="#606060" />
  </vbox>
</window>
`

const SessionTemplate = `
<window id="session" title="Session" padding="1">
  <vbox gap="0">
    <tabbar id="tabs" tabs="${tabs}" active="${activeTab}" />
    <spacer height="1" />
    <list id="messages" items="${messages}" />
    <spacer flex="1" />
    <input id="input" placeholder="Type your message..." width="100%" />
  </vbox>
</window>
`

const SetupTemplate = `
<window id="setup" title="Setup" center="true" padding="2">
  <vbox gap="1">
    <text id="title" value="Welcome to Freecode" bold="true" />
    <spacer height="1" />
    <text id="provider-label" value="Select a provider:" />
    <selectionlist id="providers" items="${providers}" selected="${selectedProvider}" />
    <spacer flex="1" />
    <hbox gap="2">
      <spacer flex="1" />
      <button id="cancel" label="Cancel" />
      <button id="next" label="Next" primary="true" />
    </hbox>
  </vbox>
</window>
`
