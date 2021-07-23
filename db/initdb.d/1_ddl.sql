CREATE SCHEMA IF NOT EXISTS `gke-go-sample` DEFAULT CHARACTER SET utf8mb4;
USE `gke-go-sample`;

CREATE TABLE IF NOT EXISTS `administrators` (
  `id` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`)
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `regions` (
  `geocode_1` CHAR(10) NOT NULL,
  `geocode_2` CHAR(20) NOT NULL,
  `zipcode` CHAR(10) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `l_area` CHAR(10) NOT NULL,
  `l_area_name` VARCHAR(50) NOT NULL,
  `original_m_area` CHAR(10) NOT NULL,
  `original_m_area_name` VARCHAR(100) NOT NULL,
  `original_s_area` CHAR(10) NOT NULL,
  `original_s_area_name` VARCHAR(100) NOT NULL,
  `m_area` CHAR(10) NOT NULL,
  `m_area_name` VARCHAR(100) NOT NULL,
  `s_area` CHAR(10) NULL,
  `s_area_name` VARCHAR(100) NULL,
  INDEX `regions_l_area_idx` (`l_area` ASC),
  INDEX `regions_m_area_idx` (`m_area` ASC),
  INDEX `regions_s_area_idx` (`s_area` ASC)
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `cities` (
  `id` VARCHAR(255) NOT NULL,
  `pref_id` CHAR(2) NOT NULL,
  `pref_name` VARCHAR(255) NOT NULL,
  `city_id` VARCHAR(255) NOT NULL,
  `city_name` VARCHAR(255) NOT NULL,
  `city_name_kana` VARCHAR(255) NOT NULL,
  `area_name` VARCHAR(255) NOT NULL,
  `area_name_kana` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `cities_pref_id_idx` (`pref_id` ASC),
  INDEX `cities_city_id_idx` (`city_id` ASC)
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `lines` (
  `id` VARCHAR(255) NOT NULL,
  `rail_id` CHAR(7) NOT NULL,
  `station_id` CHAR(5) NOT NULL,
  `stop_order` int(11) NOT NULL,
  `rail_company_name` VARCHAR(255) NOT NULL,
  `rail_company_kana` VARCHAR(255) NOT NULL,
  `rail_company_name2` VARCHAR(255) NOT NULL,
  `rail_name1` VARCHAR(255) NOT NULL,
  `rail_name_kana1` VARCHAR(255) NOT NULL,
  `rail_name2` VARCHAR(255) NOT NULL,
  `rail_name_kana2` VARCHAR(255) NOT NULL,
  `station_name` VARCHAR(255) NOT NULL,
  `station_name_kana` VARCHAR(255) NOT NULL,
  `pref_id` CHAR(2) NOT NULL,
  `latitude` double(9,6) NOT NULL,
  `longitude` double(9,6) NOT NULL,
  `rail_kind` int(11) NOT NULL,
  `rail_kind_name` VARCHAR(10) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `lines_rail_id_idx` (`rail_id` ASC),
  INDEX `lines_station_id_idx` (`station_id` ASC),
  INDEX `fk_lines_cities_idx` (`pref_id` ASC),
  CONSTRAINT `fk_lines_cities`
    FOREIGN KEY (`pref_id`)
    REFERENCES `cities` (`pref_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `agencies` (
  `id` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `name_kana` VARCHAR(255) NOT NULL,
  `postal_code` VARCHAR(255) NOT NULL,
  `pref_id` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_agencies_cities_idx` (`pref_id` ASC),
  CONSTRAINT `fk_agencies_cities`
    FOREIGN KEY (`pref_id`)
    REFERENCES `cities` (`pref_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `agency_accounts` (
  `id` VARCHAR(255) NOT NULL,
  `v1_id` VARCHAR(255) NULL,
  `agency_id` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `name_kana` VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_agency_accounts_agencies_idx` (`agency_id` ASC),
  CONSTRAINT `fk_agency_accounts_agencies`
    FOREIGN KEY (`agency_id`)
    REFERENCES `agencies` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `companies` (
  `id` VARCHAR(255) NOT NULL,
  `agency_id` VARCHAR(255) NOT NULL,
  `status` int(11) NOT NULL,
  `rank_type` int(11) NOT NULL,
  `rank` int(11) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `name_kana` VARCHAR(255) NOT NULL,
  `postal_code` VARCHAR(255) NOT NULL,
  `pref_id` CHAR(2) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `building` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_companies_agencies_idx` (`agency_id` ASC),
  CONSTRAINT `fk_companies_agencies`
    FOREIGN KEY (`agency_id`)
    REFERENCES `agencies` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_companies_cities_idx` (`pref_id` ASC),
  CONSTRAINT `fk_companies_cities`
    FOREIGN KEY (`pref_id`)
    REFERENCES `cities` (`pref_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `departments` (
  `id` VARCHAR(255) NOT NULL,
  `agency_id` VARCHAR(255) NOT NULL,
  `company_id` VARCHAR(255) NOT NULL,
  `sales_id` VARCHAR(255) NOT NULL,
  `status` int(11) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `business_condition` int(11) NOT NULL,
  `postal_code` VARCHAR(255) NOT NULL,
  `pref_id` CHAR(2) NOT NULL,
  `city_id` CHAR(5) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `building` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(255) NOT NULL,
  `m_area_id` CHAR(10) NULL,
  `s_area_id` CHAR(10) NULL,
  `latitude` double(9,6) NOT NULL,
  `longitude` double(9,6) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_departments_agencies_idx` (`agency_id` ASC),
  CONSTRAINT `fk_departments_agencies`
    FOREIGN KEY (`agency_id`)
    REFERENCES `agencies` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_departments_companies_idx` (`company_id` ASC),
  CONSTRAINT `fk_departments_companies`
    FOREIGN KEY (`company_id`)
    REFERENCES `companies` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_departments_agency_accounts_idx` (`sales_id` ASC),
  CONSTRAINT `fk_departments_agency_accounts`
    FOREIGN KEY (`sales_id`)
    REFERENCES `agency_accounts` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_departments_cities_idx` (`pref_id` ASC),
  CONSTRAINT `fk_departments_cities`
    FOREIGN KEY (`pref_id`)
    REFERENCES `cities` (`pref_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_departments_cities_idx_2` (`city_id` ASC),
  CONSTRAINT `fk_departments_cities_2`
    FOREIGN KEY (`city_id`)
    REFERENCES `cities` (`city_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_departments_regions_idx` (`m_area_id` ASC),
  CONSTRAINT `fk_departments_regions`
    FOREIGN KEY (`m_area_id`)
    REFERENCES `regions` (`m_area`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_departments_regions_idx_2` (`s_area_id` ASC),
  CONSTRAINT `fk_departments_regions_2`
    FOREIGN KEY (`s_area_id`)
    REFERENCES `regions` (`s_area`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `department_stations` (
  `id` VARCHAR(255) NOT NULL,
  `department_id` VARCHAR(255) NOT NULL,
  `line_id` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_department_stations_departments_idx` (`department_id` ASC),
  CONSTRAINT `fk_department_stations_departments`
    FOREIGN KEY (`department_id`)
    REFERENCES `departments` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_department_stations_lines_idx` (`line_id` ASC),
  CONSTRAINT `fk_department_stations_lines`
    FOREIGN KEY (`line_id`)
    REFERENCES `lines` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `department_images` (
  `id` VARCHAR(255) NOT NULL,
  `department_id` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_department_images_departments_idx` (`department_id` ASC),
  CONSTRAINT `fk_department_images_departments`
    FOREIGN KEY (`department_id`)
    REFERENCES `departments` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `works` (
  `id` VARCHAR(255) NOT NULL,
  `department_id` VARCHAR(255) NOT NULL,
  `status` int(11) NOT NULL,
  `work_type` int(11) NOT NULL,
  `job_code` int(11) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `date_from` DATETIME NOT NULL,
  `date_to` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_works_departments_idx` (`department_id` ASC),
  CONSTRAINT `fk_works_departments`
    FOREIGN KEY (`department_id`)
    REFERENCES `departments` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `work_images` (
  `id` VARCHAR(255) NOT NULL,
  `work_id` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  `view_order` int(11) NOT NULL,
  `comment` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_work_images_works_idx` (`work_id` ASC),
  CONSTRAINT `fk_work_images_works`
    FOREIGN KEY (`work_id`)
    REFERENCES `works` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `work_movies` (
  `id` VARCHAR(255) NOT NULL,
  `work_id` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_work_movies_works_idx` (`work_id` ASC),
  CONSTRAINT `fk_work_movies_works`
    FOREIGN KEY (`work_id`)
    REFERENCES `works` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `work_merits` (
  `id` VARCHAR(255) NOT NULL,
  `work_id` VARCHAR(255) NOT NULL,
  `value` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_work_merits_works_idx` (`work_id` ASC),
  CONSTRAINT `fk_work_merits_works`
    FOREIGN KEY (`work_id`)
    REFERENCES `works` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `main_contracts` (
  `id` VARCHAR(255) NOT NULL,
  `department_id` VARCHAR(255) NOT NULL,
  `status` int(11) NOT NULL,
  `plan` int(11) NOT NULL,
  `date_from` DATETIME NOT NULL,
  `date_to` DATETIME NOT NULL,
  `price` int(11) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_main_contracts_departments_idx` (`department_id` ASC),
  CONSTRAINT `fk_main_contracts_departments`
    FOREIGN KEY (`department_id`)
    REFERENCES `departments` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `work_active_plans` (
  `work_id` VARCHAR(255) NOT NULL,
  `main_contract_id` VARCHAR(255) NULL,
  `published_order` int(11) NOT NULL,
  PRIMARY KEY (`work_id`),
  INDEX `fk_work_active_plans_works_idx` (`work_id` ASC),
  CONSTRAINT `fk_work_active_plans_works`
    FOREIGN KEY (`work_id`)
    REFERENCES `works` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_work_active_plans_main_contracts_idx` (`main_contract_id` ASC),
  CONSTRAINT `fk_work_active_plans_main_contracts`
    FOREIGN KEY (`main_contract_id`)
    REFERENCES `main_contracts` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `entries` (
  `id` VARCHAR(255) NOT NULL,
  `department_id` VARCHAR(255) NOT NULL,
  `work_id` VARCHAR(255) NOT NULL,
  `full_name` VARCHAR(255) NOT NULL,
  `full_name_kana` VARCHAR(255) NOT NULL,
  `birthdate` DATE NOT NULL,
  `gender` int(11) NOT NULL,
  `phone_number` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `question` VARCHAR(255) NOT NULL,
  `category` VARCHAR(255) NULL,
  `pref_id` VARCHAR(255) NULL,
  `preferred_contact_method` int(11) NULL,
  `preferred_contact_time` VARCHAR(255) NULL,
  `status` int(11) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_entries_departments_idx` (`department_id` ASC),
  CONSTRAINT `fk_entries_departments`
    FOREIGN KEY (`department_id`)
    REFERENCES `departments` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_entries_works_idx` (`work_id` ASC),
  CONSTRAINT `fk_entries_works`
    FOREIGN KEY (`work_id`)
    REFERENCES `works` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `usage_statements` (
  `id` VARCHAR(255) NOT NULL,
  `department_id` VARCHAR(255) NOT NULL,
  `main_contract_id` VARCHAR(255) NULL,
  `price` int(11) NOT NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_usage_statements_departments_idx` (`department_id` ASC),
  CONSTRAINT `fk_usage_statements_departments`
    FOREIGN KEY (`department_id`)
    REFERENCES `departments` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  INDEX `fk_usage_statements_main_contracts_idx` (`main_contract_id` ASC),
  CONSTRAINT `fk_usage_statements_main_contracts`
    FOREIGN KEY (`main_contract_id`)
    REFERENCES `main_contracts` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;