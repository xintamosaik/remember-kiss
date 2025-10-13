import sys
import os

todos = []

def run(message = "starting system.."):
    os.system('cls' if os.name == 'nt' else 'clear')
    print(message)
    print("\nTodo List:")
    if not todos:
        print("  (No items)")
    for i, todo in enumerate(todos, 1):
        print(f"  {i}. {todo}")
    
    print("\nOptions:")
    print("  a - Add new todo")
    print("  q - Quit")
    
    choice = input("Select item number to edit, or option: ").strip()
    if choice == 'q':
        print("Goodbye!")
        sys.exit(0)

    if choice == 'a':
        text = input("Enter new todo: ").strip()
        if text:
            todos.append(text)
            run("Todo added.")
        else:
            run("Empty todo not added.")
    
    if choice.isdigit() == False:
        run("Unknown option.")

    idx = int(choice) - 1
    no_match_for_number = idx < len(todos)
    if no_match_for_number == False:
        run("Invalid number.")

    print(f"Current: {todos[idx]}")
    new_text = input("Enter new text (leave empty to cancel): ").strip()
    if new_text:
        todos[idx] = new_text
        run("Todo updated.")
    else:
        run("Edit cancelled.")
    

if __name__ == "__main__":
    run()