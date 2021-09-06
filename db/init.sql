

CREATE TABLE `job` (
                       `id` bigint NOT NULL AUTO_INCREMENT,
                       `name` varchar(30) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'unique name',
                       `description` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'job description',
                       `cron` varchar(20) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'cron expression',
                       `next_time` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'job run next time',
                       `status` int DEFAULT NULL COMMENT 'job status',
                       `concurrent` int NOT NULL DEFAULT '0' COMMENT '是否可以并发执行',
                       PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



CREATE TABLE `job_lock` (
                            `lock_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                            `lock_val` bigint DEFAULT NULL,
                            PRIMARY KEY (`lock_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



CREATE TABLE `job_log` (
                           `id` bigint NOT NULL AUTO_INCREMENT,
                           `job_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                           `start_time` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
                           `end_time` varchar(20) COLLATE utf8mb4_general_ci DEFAULT NULL,
                           `status` int NOT NULL,
                           `error_msg` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO job_lock (lock_name,lock_val) VALUES ('job_execute',0);