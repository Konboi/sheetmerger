# sheetmerger

Merge sheets managed with special rules

## Useage

```sh
 sheetmerger -c <CONFIG YML FILE>  -base <BASE SPREAD SHEET ID> -diff <DIFF SPREAD SHEET ID > -name <> -name <SHEET NAME1> -name <SHEET NAME2>
 ```

## Config

```yml
client:
  email: {{ must_env "SPREADSHEET_SERVICE_EMAIL" }}
  private_key_id: {{ must_env "SPREADSHEET_PRIVATE_KEY_ID" }}
  private_key: "{{ must_env `SPREADSHEET_PRIVATE_KEY` }}"
base_sheet_name: {{ env "BASE_SHEET_NAME" }}
sheet_index_column: {{ env "INDEX_COLUMN" }}
backup_folder_id: {{ must_env "BACKUP_FOLDER_ID"}}
```