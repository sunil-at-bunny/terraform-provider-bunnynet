---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bunny_video_language Data Source - terraform-provider-bunny"
subcategory: ""
description: |-
  Video Language
---

# bunny_video_language (Data Source)

Video Language

## Example Usage

```terraform
data "bunny_video_language" "en" {
  code = "en"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `code` (String)

### Read-Only

- `name` (String)
- `support_player_translation` (Boolean)
- `support_transcribing` (Boolean)
- `transcribing_accuracy` (Number)