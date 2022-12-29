import os


class FileManager:
    @classmethod
    def dir_exists(self, path: str) -> bool:
        return os.path.exists(path)

    @classmethod
    def get_or_create_dir(self, path: str) -> str:
        if not self.dir_exists(path):
            os.mkdir(path)
        return path

    @classmethod
    def list_file_names(self, dir_path: str) -> list[str]:
        return os.listdir(self.get_or_create_dir(dir_path))

    @classmethod
    def count_files(self, dir_path: str) -> int:
        return len(self.list_file_names(dir_path))
