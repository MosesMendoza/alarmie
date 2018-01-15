**Alarmie**

A slack bot for ops event/incident management.

**Status**

Active development - does not work yet

**Premise**

The hardest part of incident management (aside from fixing the incident) is communication. Alarmie serves as a central point of communication about the status of an ongoing incident.

**Interface**

Communicate with alarmie directly as a bot in slack.
```
@alarmie
  incident
    create => #number
    close #number => ok
    reopen #number => ok

  status #number
    update <message> => ok
    get => <message>
    log => <messages>

  assignee #number
    get => <message>
    update <message> => ok

  history
    get => <#numbers>
    get #number => <message>
```

**License**
MIT
See [LICENSE](LICENSE)

