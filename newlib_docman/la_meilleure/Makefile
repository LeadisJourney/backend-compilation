CC          =	gcc -std=c11

MKDIR       =   mkdir -p

RM          =   rm -rf

CCFLAGS     =	-Wfatal-errors -Werror -Wall -Wextra -pedantic

LDFLAGS     =	-L./lib/ -Wl,-rpath=./lib/

LDLIBS      =	-llaMeilleure

PRG         =	res

INCS_DIR    =   ./inc/

LIBS_DIR    =   ./lib/

OBJS_DIR    =   ./obj/

SRCS_DIR    =   ./src/

SRCS	    =	$(SRCS_DIR)la_meilleure.c

OBJS        :=  $(SRCS:$(SRCS_DIR)%.c=$(OBJS_DIR)%.o)

DEPS        :=  $(OBJS:$(OBJS_DIR)%.o=$(OBJS_DIR)%.d)

#---RULES---
all: objDir $(PRG)

objDir:
	@$(MKDIR) $(OBJS_DIR)

-include $(DEPS)

$(PRG): $(OBJS)
	@$(CC) $(OBJS) -o $@ $(LDFLAGS) $(LDLIBS)

$(OBJS_DIR)%.o: $(SRCS_DIR)%.c
	@$(CC) -c $< -I$(INCS_DIR) -I$(TPLS_DIR) -I$(LIBS_DIR) $(CCFLAGS) -o $@
	@$(CC) -M $< -I$(INCS_DIR) -I$(TPLS_DIR) -I$(LIBS_DIR) -MT $@ -MF $(OBJS_DIR)$*.d

clean:
	@$(RM) $(OBJS) $(DEPS) $(OBJS_DIR)

distclean: clobber
realclean: clobber
clobber: clean
	@$(RM) $(PRG)

.PHONY: all clean distclean realclean clobber
