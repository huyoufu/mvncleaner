/*
    索引表
    字段:
        id id
        name 文件名字
        file_type 文件类型
        created 创建时间
        modify 修改时间
 */
CREATE TABLE IF NOT EXISTS "mc_index"
(
    "id"   INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" VARCHAR(255) NULL ,
    "file_type" INTEGER NULL ,
    "created"  TIMESTAMP default (datetime('now', 'localtime')),
    "modify"  TIMESTAMP default (datetime('now', 'localtime'))
);
/**
  * 给文件名字段加索引
 */
create index mc_index_name on mc_index (name);











