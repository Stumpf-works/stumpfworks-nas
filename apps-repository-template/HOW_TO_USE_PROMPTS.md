# 📖 How to Use the Prompts in This Repository

This repository contains specialized prompts for working with Claude Code effectively.

---

## 🎯 Available Prompts

### 1. **QUICK_START_PROMPT.txt** ⚡
**When to use**: At the start of EVERY Claude Code session in this repository

**What it does**:
- Gives Claude immediate context about where he is
- Explains basic structure and rules
- Lists key commands and workflows

**How to use**:
```
Copy-paste the entire content of QUICK_START_PROMPT.txt into your first message:

[Paste QUICK_START_PROMPT.txt content]

Then add your actual request:

I want to add a MinIO plugin for S3-compatible storage.
```

**Expected result**: Claude understands the repository and can start working immediately.

---

### 2. **SESSION_PROMPT.md** 📋
**When to use**: When you need more detailed context or reference during development

**What it does**:
- Complete repository structure explanation
- All workflows in detail
- Security guidelines
- Common scenarios with solutions
- Quick reference commands

**How to use**:
```
Reference it with @ in Claude Code:

@SESSION_PROMPT.md I'm getting validation errors, what should I check?
```

**Or**: Open it in your editor and read specific sections when needed.

---

### 3. **CLAUDE_CODE_MASTER_PROMPT.md** 🏗️
**When to use**: When developing plugins that integrate with StumpfWorks NAS

**What it does**:
- Explains StumpfWorks NAS architecture completely
- System Library APIs
- Plugin development patterns
- Backend & frontend integration
- Database models
- Best practices

**How to use**:
```
Reference when you need NAS-specific context:

@CLAUDE_CODE_MASTER_PROMPT.md How do I access the ZFS API from my plugin?

@CLAUDE_CODE_MASTER_PROMPT.md What's the correct pattern for database models?
```

**This is THE authoritative source** for StumpfWorks NAS plugin development!

---

## 🎬 Example Session Flows

### Scenario 1: Starting Fresh - New Plugin

```
Step 1: Copy QUICK_START_PROMPT.txt content into first message
Step 2: Add your request
```

**Example**:
```
[QUICK_START_PROMPT.txt content]

I want to create a Nextcloud plugin that integrates with StumpfWorks NAS
user management. Help me set it up.
```

**Claude will**:
- Understand repository context ✓
- Know to use templates ✓
- Follow proper plugin structure ✓
- Reference master prompt for NAS APIs ✓

---

### Scenario 2: Continuing Work - Existing Plugin

```
Step 1: Brief context
Step 2: Reference specific prompt if needed
```

**Example**:
```
I'm working on the MinIO plugin in plugins/minio/. I need to update it to
version 1.1.0 with a new feature.

@SESSION_PROMPT.md section on updating existing plugins
```

**Claude will**:
- Know the update workflow ✓
- Update version correctly ✓
- Update CHANGELOG ✓
- Create proper release ✓

---

### Scenario 3: Integration Question - NAS APIs

```
Reference MASTER_PROMPT for architecture
```

**Example**:
```
@CLAUDE_CODE_MASTER_PROMPT.md

My plugin needs to create ZFS datasets. How do I access the System Library
from my plugin?
```

**Claude will**:
- Know the System Library architecture ✓
- Explain API access patterns ✓
- Provide correct code examples ✓

---

### Scenario 4: Validation Error - Troubleshooting

```
Reference SESSION_PROMPT for workflows
```

**Example**:
```
@SESSION_PROMPT.md

I ran validate-plugins.py and got errors. What do I need to fix?

Error: Missing required field: category
```

**Claude will**:
- Know validation requirements ✓
- Explain how to fix ✓
- Show correct plugin.json format ✓

---

## 🎯 Best Practices

### ✅ DO:

1. **Always start with QUICK_START_PROMPT.txt**
   - Copy entire content
   - Paste at start of first message
   - Add your request after

2. **Use @ references for specific help**
   - `@SESSION_PROMPT.md` for repository workflows
   - `@CLAUDE_CODE_MASTER_PROMPT.md` for NAS integration

3. **Be specific about what you're working on**
   ```
   Good: "I'm updating plugins/asterisk-voip/ to v1.1.0"
   Bad:  "Update the plugin"
   ```

4. **Reference sections when possible**
   ```
   @SESSION_PROMPT.md section on "Adding a New Plugin"
   ```

### ❌ DON'T:

1. **Don't assume Claude knows where you are**
   - Always provide context at session start
   - Use QUICK_START_PROMPT.txt

2. **Don't mix contexts**
   - This repo is for DISTRIBUTION
   - For NAS core development, use different prompts

3. **Don't skip validation**
   - Always run scripts before asking for review

---

## 📊 Prompt Hierarchy

```
┌─────────────────────────────────────────┐
│   QUICK_START_PROMPT.txt                │  ← Start here!
│   (Quick context, basic info)           │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│   SESSION_PROMPT.md                     │  ← Reference during work
│   (Detailed workflows, commands)        │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│   CLAUDE_CODE_MASTER_PROMPT.md          │  ← Reference for NAS integration
│   (StumpfWorks NAS architecture)        │
└─────────────────────────────────────────┘
```

---

## 🔄 Session Lifecycle

### Beginning of Session:
1. Open Claude Code
2. Copy QUICK_START_PROMPT.txt
3. Paste + add your task
4. Claude is now fully contextualized ✓

### During Session:
- Reference SESSION_PROMPT.md for workflows
- Reference MASTER_PROMPT for NAS APIs
- Claude maintains context throughout

### End of Session:
- Validate: `python3 scripts/validate-plugins.py`
- Test locally if possible
- Commit with proper message
- Tag if releasing

### Next Session:
- Start over with QUICK_START_PROMPT.txt
- Claude doesn't remember previous sessions
- Always re-establish context

---

## 💡 Pro Tips

1. **Keep prompts open in editor tabs**
   - Easy to reference
   - Quick copy-paste
   - Search within prompts

2. **Use specific sections**
   ```
   Not: "How do I add a plugin?"
   Better: "@SESSION_PROMPT.md workflow 1️⃣ Adding a New Plugin"
   ```

3. **Combine prompts when needed**
   ```
   I'm creating a new plugin that uses ZFS.

   @SESSION_PROMPT.md (for repository workflow)
   @CLAUDE_CODE_MASTER_PROMPT.md (for ZFS API access)
   ```

4. **Save common tasks as snippets**
   - Your editor can save frequently used prompts
   - Create shortcuts for common scenarios

---

## 🎓 Learning Curve

**First Time**:
- Read all three prompts completely
- Understand repository structure
- Try adding a simple plugin

**Regular Use**:
- Quick start prompt at session start
- Reference others as needed
- Workflows become natural

**Advanced**:
- Know which prompt has what info
- Jump directly to relevant sections
- Fast, efficient development

---

## 📱 Quick Reference Card

Save this for easy access:

```
┌──────────────────────────────────────────────┐
│ CLAUDE CODE - STUMPFWORKS NAS APPS          │
├──────────────────────────────────────────────┤
│                                              │
│ START SESSION:                               │
│   Copy QUICK_START_PROMPT.txt               │
│                                              │
│ REPOSITORY HELP:                             │
│   @SESSION_PROMPT.md                         │
│                                              │
│ NAS INTEGRATION:                             │
│   @CLAUDE_CODE_MASTER_PROMPT.md             │
│                                              │
│ VALIDATE:                                    │
│   python3 scripts/validate-plugins.py        │
│                                              │
│ REGISTRY:                                    │
│   python3 scripts/generate-registry.py       │
│                                              │
└──────────────────────────────────────────────┘
```

---

## ❓ FAQ

**Q: Do I need to use prompts every time?**
A: Yes! Claude Code sessions are stateless. Each new session needs context.

**Q: Can I just say "continue working on X"?**
A: No - Claude won't remember. Always re-establish context with QUICK_START_PROMPT.txt

**Q: Which prompt for what?**
- QUICK_START → Start of session (always!)
- SESSION_PROMPT → Repository workflows
- MASTER_PROMPT → NAS architecture & APIs

**Q: Can I modify the prompts?**
A: Yes, but keep them up-to-date with repository changes. They're documentation too!

**Q: What if Claude seems confused?**
A: Likely missing context. Provide relevant prompt again.

---

## 🎉 You're Ready!

You now know how to effectively work with Claude Code in this repository.

**Remember**:
- Context is EVERYTHING
- Start with QUICK_START_PROMPT.txt
- Reference others as needed
- Validate before committing

**Happy Plugin Development! 🚀**
