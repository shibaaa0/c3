import curses
import os
import subprocess

def init_colors():
    curses.start_color()
    curses.init_pair(1, curses.COLOR_YELLOW, curses.COLOR_BLACK)    # Đường dẫn
    curses.init_pair(2, curses.COLOR_BLUE, curses.COLOR_BLACK)      # Thư mục
    curses.init_pair(3, curses.COLOR_WHITE, curses.COLOR_BLACK)     # File, Thanh tìm kiếm
    curses.init_pair(4, curses.COLOR_WHITE, curses.COLOR_CYAN)      # Mục được chọn

def get_search_input(stdscr, height, width):
    stdscr.clear()
    prompt = "Search: "
    stdscr.addstr(height - 1, 0, prompt, curses.color_pair(3) | curses.A_BOLD)
    curses.echo()
    search_term = stdscr.getstr(height - 1, len(prompt), width - len(prompt) - 1).decode('utf-8')
    curses.noecho()
    stdscr.clear()
    return search_term.strip()

def search_files(current_dir, search_term):
    all_items = os.listdir(current_dir)
    return [item for item in all_items if search_term.lower() in item.lower()]

def main(stdscr):
    curses.curs_set(0)  
    if curses.has_colors():
        init_colors()
    
    current_dir = os.getcwd()
    selected = 0
    show_hidden = False
    search_mode = False
    search_results = []
    
    while True:
        if search_mode and search_results:
            items = search_results
        else:
            if show_hidden:
                items = os.listdir(current_dir)
            else:
                items = [item for item in os.listdir(current_dir) if not item.startswith('.')]
            items.sort()
        stdscr.clear()
        
        height, width = stdscr.getmaxyx()
        
        path_display = f">>{current_dir}"[:width-1]
        if curses.has_colors():
            stdscr.addstr(0, 0, path_display, curses.color_pair(1) | curses.A_BOLD)
        else:
            stdscr.addstr(0, 0, path_display)
        
        for i, item in enumerate(items):
            if i >= height - 1:  
                break
            display_text = item[:width-1]
            attrs = 0
            if curses.has_colors():
                if i == selected:  
                    attrs = curses.color_pair(4) | curses.A_BOLD
                elif os.path.isdir(os.path.join(current_dir, item)):
                    attrs = curses.color_pair(2) | curses.A_BOLD
                else:
                    attrs = curses.color_pair(3) | curses.A_BOLD
            stdscr.addstr(i + 1, 0, display_text, attrs)  
        stdscr.refresh()

        key = stdscr.getch()
        
        # Xử lý phím điều hướng
        if key == curses.KEY_UP and selected > 0:
            selected -= 1
        elif key == curses.KEY_DOWN and selected < len(items) - 1:
            selected += 1
        elif key in {ord('i'), curses.KEY_ENTER, 10, 13}:  # Mở thư mục hoặc file
            selected_path = os.path.join(current_dir, items[selected])
            if os.path.isdir(selected_path):
                current_dir = os.path.abspath(selected_path)
                selected = 0
                search_mode = False
                stdscr.clear()
            elif os.path.isfile(selected_path):
                curses.endwin()
                subprocess.run(['micro', selected_path])
                stdscr.clear()
                stdscr.refresh()
        
        
        elif key == ord('f'):  # Tìm kiếm 
            search_term = get_search_input(stdscr, height, width)
            if search_term:
                search_results = search_files(current_dir, search_term)
                if search_results:
                    search_mode = True
                    selected = 0
                else:
                    search_mode = False
            stdscr.clear()
        
     
        elif key == ord('s'):  # Show hidden file
            show_hidden = not show_hidden
            search_mode = False  
            selected = 0
            stdscr.clear()
               
        elif key == ord('o'):  # Quay lại thư mục cha 
            new_dir = os.path.dirname(current_dir)
            if new_dir != current_dir:
                current_dir = new_dir
                selected = 0
                search_mode = False
                stdscr.clear()

if __name__ == "__main__":
    try:
        curses.wrapper(main)
    except KeyboardInterrupt:
        pass
    except Exception:
        pass